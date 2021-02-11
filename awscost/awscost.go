package awscost

import (
	"log"
	"time"

	"../config"
	"../utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"gopkg.in/mgo.v2/bson"
)

type Cost struct {
	Cost string    // `json:"cost" bson:"cost"`
	Date time.Time // `json:"date" bson:"date"`
}

// Sends a query to AWS Cost Explorer if the latest record
// was created at least 96 hours (4 days) ago
// or if the data is not present yet it will proceed
// and updates homeassistant sensor.aws_monthly_cost value
func Update() {
	cost := getRecentCost()
	now := time.Now()
	if !(now.Sub(cost.Date).Hours() > 96) && cost.Cost != "" {
		updateHassio(cost.Cost)
		return
	}
	currentMonth := time.Now().Format("2006-01")
	nextMonth := time.Now().AddDate(0, 1, 0).Format("2006-01")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc := costexplorer.New(sess)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := svc.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		TimePeriod: &costexplorer.DateInterval{
			Start: aws.String(currentMonth + "-01"),
			End:   aws.String(nextMonth + "-01"),
		},
		Granularity: aws.String("MONTHLY"),
		Metrics: aws.StringSlice([]string{
			"BlendedCost",
		}),
	})
	if err != nil {
		log.Fatalf("[AWSCOST] Unable to generate report, %v/n", err)
	}
	currentCost := *result.ResultsByTime[0].Total["BlendedCost"].Amount
	insertCost(currentCost)
	updateHassio(currentCost)

	log.Println("[AWSCOST] Cost Report:", currentCost)
}

func getRecentCost() Cost {
	cost := Cost{}
	err := config.AWS.Find(nil).Sort("-date").Limit(1).One(&cost)
	if err != nil {
		log.Fatalln(err)
	}

	return cost
}

func insertCost(cost string) {
	err := config.AWS.Insert(bson.M{"cost": cost, "date": time.Now()})
	if err != nil {
		log.Println(err)
	}
}

func updateHassio(cost string) {
	sensor := utils.Sensor{
		State: cost,
		Attributes: utils.Attributes{
			Friendly_name:       "AWS monthly cost",
			Unit_of_measurement: "$",
			Icon:                "mdi:currency-usd",
		},
	}
	utils.SetState("sensor.aws_monthly_cost", sensor)
}
