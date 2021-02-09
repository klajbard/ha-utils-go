package awscost

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

func Update() {
	currentMonth := time.Now().Format("2006-01")
	nextMonth := time.Now().AddDate(0, 1, 0).Format("2006-01")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)
	svc := costexplorer.New(sess)

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

	log.Println("[AWSCOST] Cost Report:", *result.ResultsByTime[0].Total["BlendedCost"].Amount)
}
