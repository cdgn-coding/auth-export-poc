# Auth0 Users Sync Job POC

This project synchronizes users from an Auth0 database 
with a PostgreSQL with the goal of being able to perform queries
that are out of the scope of the Auth0's Management API.

## Requirements

* Extract users from Auth0
* Enrich data by adding the linked organizations
* Load users to PostgreSQL