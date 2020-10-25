# zipcode2address

Searching japanese address by zipcode.
Using serverless framework and Golang.

# feature

## import

Importing japanese address data to DynamoDB.

## search

Searching address by zipcode.

# deploy

1. Prepare environments

```
mkdir env
touch env/dev.yml
touch env/prod.yml
```

2. deploy

```
sh deploy.sh
```
