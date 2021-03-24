package main

import (
	"context"
	"fmt"
	"time"
  "net/http"
  "os"
  "crypto/tls"
  "database/sql"

  _ "github.com/lib/pq"

  "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

// Test variables
// Postgresql Operator
const (
  host = "acid-minimal-cluster-repl"
  port = 5432
  user = "postgres"
  password = "EgOQPibRIVIvgOEnQIJis4qngZ2WWBSRJP0sBGc4R4eKIxWSXL1sNn58Q70exWS8"
  dbname = "postgres"
)

func curlElastic (){
// curl -u "elastic:$PASSWORD" -k "https://quickstart-es-http:9200"
// TODO: This is insecure; use only in dev environments.
  tr := &http.Transport{
	  TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
  }
  client := &http.Client{Transport: tr}

  req, err := http.NewRequest("GET", "https://quickstart-es-http:9200", nil)
  if err != nil {
    fmt.Println("Error creating request  to Elasticsearch")
    fmt.Print(err)
  }
  req.SetBasicAuth("elastic", os.ExpandEnv("$PASSWORD"))

  resp, err := client.Do(req)
  if err != nil {
    fmt.Println("Error connection to Elasticsearch")
    fmt.Print(err)
  }
  defer resp.Body.Close()

  fmt.Println(resp)
}

func connectPsql (){

  // connection string
  psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

  // open database
  db, err := sql.Open("postgres", psqlconn)
  CheckError(err)

  // close database
  defer db.Close()

  // check db
  err = db.Ping()
  CheckError(err)

  fmt.Println("Connected to Database!")
}

func CheckError(err error){
 if err != nil {
  panic(err)
 }
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace, in this case default
		pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster in the default namespace\n", len(pods.Items))

   // Make call to Elastic
    curlElastic()
   // Connect to DB
    connectPsql()



   // TODO: Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
		_, err = clientset.CoreV1().Pods("default").Get(context.TODO(), "misc-go", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("Pod misc-go not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		} else {
			fmt.Printf("Found misc-go  pod in default namespace\n")
		}

		time.Sleep(10 * time.Second)
	}
}

