#!/bin/bash
curl -X PUT "http://ec2-52-34-228-148.us-west-2.compute.amazonaws.com:8080/v2/apps/router" -d @"router.json" -H "Content-type: application/json"
