# Country Information Service

This is a RESTful web service that provides information about countries, including population data and city lists.

## Endpoints

### 1. `/countryinfo/v1/info/{country_code}`
- Returns general information about a country.
- Example: `http://localhost:8080/countryinfo/v1/info/no?limit=5`

### 2. `/countryinfo/v1/population/{country_code}`
- Returns historical population data for a country.
- Example: `http://localhost:8080/countryinfo/v1/population/no?limit=2010-2015`

### 3. `/countryinfo/v1/status/`
- Returns the status of the service and external APIs.
- Example: `http://localhost:8080/countryinfo/v1/status/`

## Deployment
The service is deployed on Render: [Render URL](#)