
This repository contains the source code for a take-home exam for Fetch Rewards. The assignment involves building a server using Go and implementing two endpoints:

The /receipts/process endpoint accepts a Receipt JSON object via a POST request. The server parses the receipt object and calculates the points based on a predefined set of rules. The endpoint returns a JSON object containing the ID of the submitted receipt, which can be used for later reference.

The /receipts/{id}/points endpoint retrieves the points associated with a previously submitted receipt. The ID of the receipt is passed as a parameter in the endpoint's path. The endpoint returns a JSON object containing the total points earned by the receipt.

The server handles various error cases, such as missing attributes in the JSON body of the POST request, errors in parsing date and time attributes, and internal server errors.

It's worth noting that the code related to this assignment can be found in the file fetch-points.go