package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

func main() {

	/*
		    -sql, -out 必须设置
			-sql "select * from app01_user_login limit 1" 
			-out s3://aws-hands-on-athena/query/
	*/

	region := flag.String("region", "us-east-1", "region default us-east-1")
	db := flag.String("db", "default", "ahtena database is 'default' ")
	sql := flag.String("sql", "", "sql, etc. select * from app01_user_login limit 10")
	out := flag.String("out", "", "outputLocation,etc. s3://path/to/query/bucket/")

	flag.Parse()

	if *sql == "" {
		fmt.Println(fmt.Errorf("SQL is required, -sql"))
		flag.Usage()
		os.Exit(-1)
	}

	if *out == "" {
		fmt.Println(fmt.Errorf("OutputLocation is required, -out"))
		flag.Usage()
		os.Exit(-1)
	}

	awscfg := &aws.Config{}
	// 指定region
	awscfg.WithRegion(*region)
	// 创建session实例
	sess := session.Must(session.NewSession(awscfg))

	// 创建athena服务实例
	svc := athena.New(sess, aws.NewConfig().WithRegion(*region))
	var s athena.StartQueryExecutionInput
	s.SetQueryString(*sql)

	var q athena.QueryExecutionContext
	q.SetDatabase(*db)
	s.SetQueryExecutionContext(&q)

	var r athena.ResultConfiguration
	r.SetOutputLocation(*out)
	s.SetResultConfiguration(&r)

	result, err := svc.StartQueryExecution(&s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("StartQueryExecution result:")
	fmt.Println(result.GoString())

	var qri athena.GetQueryExecutionInput
	qri.SetQueryExecutionId(*result.QueryExecutionId)

	var qrop *athena.GetQueryExecutionOutput
	duration := time.Duration(2) * time.Second // Pause for 2 seconds

	for {
		qrop, err = svc.GetQueryExecution(&qri)
		if err != nil {
			fmt.Println(err)
			return
		}
		if *qrop.QueryExecution.Status.State != "RUNNING" && *qrop.QueryExecution.Status.State != "QUEUED" {
			break
		}
		fmt.Printf("Waiting... %s \n", *qrop.QueryExecution.Status.State)
		time.Sleep(duration)

	}
	if *qrop.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(*result.QueryExecutionId)

		op, err := svc.GetQueryResults(&ip)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%+v", op)
	} else {
		fmt.Println(*qrop.QueryExecution.Status.State)

	}
}

