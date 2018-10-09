package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func getParametersByPath() {

}

func getParameter() {
	if len(os.Args) != 2 {
		exitErrorf("SSM Key name required\nUsage: %s key_name",
			os.Args[0])
	}

	key_name := os.Args[1]

	// Create S3 service client
	svc := ssm.New(session.New())

	params := &ssm.GetParameterInput{
		Name:           aws.String(key_name),
		WithDecryption: aws.Bool(true),
	}

	// Get the list of items
	resp, err := svc.GetParameter(params)

	if err != nil {
		exitErrorf("Unable to get key %q, %v", key_name, err)
	}

	fmt.Println(*resp.Parameter.Value)
}

func main() {

	var path = flag.String("path", "", "parameter path")
	var value = flag.String("name", "", "single parameter value")
	flag.Parse()

	// Create SSM service client
	svc := ssm.New(session.New())

	if *path == "" {
		params := &ssm.GetParameterInput{
			Name:           aws.String(*value),
			WithDecryption: aws.Bool(true),
		}
		resp, err := svc.GetParameter(params)
		if err != nil {
			exitErrorf("Unable to get key %q, %v", *value, err)
		}
		fmt.Println(*resp.Parameter.Value)

	} else {
		params := &ssm.GetParametersByPathInput{
			Path:           aws.String(*path),
			Recursive: 			aws.Bool(true),
			WithDecryption: aws.Bool(true),
		}
		resp, err := svc.GetParametersByPath(params)
		if err != nil {
			exitErrorf("Unable to get key %q, %v", *path, err)
		}
		for _, v := range resp.Parameters {
			the_key := strings.Split(*v.Name, "/")
			fmt.Printf("%s=%s;", the_key[len(the_key)-1], *v.Value)
		}

	}

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}