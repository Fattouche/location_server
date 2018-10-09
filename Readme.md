# Location_server

## Goal

The goal of this program is to accept a search string of format "/search?searchTerm=camera&lat=51.948&lng=0.172943"
and return a list of top 20 items that are related to that product and are also close in proximity.

## Usage

__you will need gcc 64bit because the sqlite package is based off cgo__

1. `go get github.com/mattn/go-sqlite3` 

2. `go get github.com/jbrukh/bayesian`

3. `go run server.go productClassifier.go geography.go`

4. Preform a GET request to `localhost/search?searchTerm=camera&lat=51.948&lng=0.172943` (I use postman for http requests, but whatever works for you!)

## Solution

This implementation uses the naive bayes machine learning classifier to sort products into specific categories. This was done by first determining a dictionary of words that are commonly seen with already known classified products. For example the word "Canon" is commonly associated with cameras and would be a good indication of the product fitting into the photography category. After the dictionaries were made, the database was trained on this model and each of those items were sorted into their corresponding categories. The next step was to sort the database items based off the parameters used in the original search. This was done in two steps, the first step was to check if the string in the database contained the search parameter, the second was to sort based off distance. This sorting was chosen because I believe that it is more important to have the correct product that is far than an incorrect product that is nearby. Lastly, the top 20 items were return to the user.

## Code Design

The code is split up into three different concerns; Server, productClassification and geographical math. The server is responsible for handling requests and parsing the user information. The productClassifier is responsible for determining what the product might be and finding the top 20 similar products. The geography file is taken from [here](https://gist.github.com/cdipaolo/d3f8db3848278b49db68) and is used to determine how close a product is to the searched product.

## Design Decisions

Golang was chosen for the language here because I have experience with it and I find it is quite intuitive when building servers like this one (It is also super fast!). I chose to classify products using the naive bayes machine learning algorithm because I find that it gives quite good results and also scales nicely once the dictionary of words you provide it is sufficiently large and accurate. github.com/jbrukh/bayesian was chosen because it is relatively well known for ML algorithms in golang and has 10 contributors. Likewise,github.com/mattn/go-sqlite3 was a no brainer choice because it is well established and maintained by over 100 contributors. I split up the code into the three sections of product classification, server and geography because I felt that it was the best way to separate concerns away from specific components. For example, the server doesn't need to know how the 20 items are getting fetched, it just wants the database (productClassifier in this case) to go and get them. In the future I would make the database queries more scalable and avoid storing all of the database items in memory(yikes I know). In addition, I would move the classification text documents into the database and add a ton more words, these were just the ones that I could think of off the top of my head.

## Example output

`localhost/search?searchTerm=canon&lat=51.948&lng=0.172943`

```
canon 400mm lense
canon 5d m3 kit
canon 24- 105mm
canon 50mm
sigma 35mm art lens 1.4 (ef / canon fit)
canon ef 70-200mm f2.8 l is ii usm lens
carl zeiss 18mm f2.8 milvus slr lens canon ze fit
canon 24-105mm f1.4 l is usm lens
canon ef 24mm f2.8 is usm lens
carl zeiss distagon 25mm f2 lens canon ze fit
canon 5dsr dlsr camera & 3 lenses and tripod
canon ef 24-70mm f2.8 l ii usm lens
canon 5dsr dlsr camera bundle including 3 lenses of your chioce
canon 5dsr digital camera
canon 17-40mm l usm f4 lens
canon 1.4 extender ii
carl zeiss 135mm f2 milvus telephoto slr lens canon ze fit
carl zeiss 50mm f1.4 milvus standard slr lens canon ze fit
canon ef 70mm - 300mm  f4 - f5.6 is usm zoom lens
canon 24-105mm 1:4 l is usm lens
canon 5d mk 1
canon ef 100-400mm lens f4.5 l is ii usm
```