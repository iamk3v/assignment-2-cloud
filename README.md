# Countries Dashboard Service

This project is an implementation of a RESTful web service built in Go for Assignment 2 of PROG2005. It enables clients to configure and retrieve dynamically populated dashboards with country data, register webhooks for notifications, and monitor service status. The service uses Firestore for persistent storage, integrates with external APIs (for country, weather, and currency data), and includes caching with periodic purging of cached data.

## Contributors
- Marius: Registrations, Testing and Stub
- Mathias: Notifications, Firebase and Testing
- Sebastian: Dashboard & Status
- Johannes: Caching, Purging and Clients

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Setup & Installation](#setup--installation)
- [Running the Application](#run-the-application)
- [API Endpoints](#api-endpoints)
- [Caching and Purging](#caching-and-purging)
- [Testing](#testing)
- [Deployment](#deployment)

## Overview

The Countries Dashboard Service allows users to:

- **Configure Dashboard Information**  
  Register, update, retrieve, and delete dashboard configurations containing country details, weather data, and currency exchange rates.

- **Retrieve Populated Dashboards**  
  Combine data from external APIs (REST Countries, Open Meteo, Currency API) based on a configuration, presenting the enriched dashboard.

- **Manage Notifications via Webhooks**  
  Register, update, retrieve, and delete webhooks that trigger notifications on events (REGISTER, CHANGE, DELETE, INVOKE).

- **Monitor Service Status**  
  Check the health of external APIs, view the number of registered webhooks, and monitor service uptime.

- **Caching**  
  Cache external API responses in Firestore and automatically purge stale cache entries using a configured TTL.

## Features

- **Registration Endpoints:**  
  Create, read, update, and delete dashboard configurations.
- **Dashboard Endpoint:**  
  Retrieve dashboards with addition of external API data.
- **Notification Endpoints:**  
  Manage webhook registrations and trigger notifications on specific events.
- **Status Endpoint:**  
  Provide system health information including external API statuses and uptime.
- **Caching:**  
  Implement caching for external API responses.
- **Testing:**  
  Comprehensive endpoint tests using Go’s `httptest` package.

### Additional Features

- **Purging of Cached Data**
- **`PATCH` method on `/notifications/` and `/registrations/`**
- **`HEAD` method on `/notifications/` and `/registrations/`**
- **Timezone Information in any time representation**


## Setup & Installation

### Prerequisites

- **Go:** Version 1.24.1 or higher.
- **Firestore:** A Google Cloud Firestore project with a service account JSON file placed at `config/service-account.json`
- **Docker (Optional):** For containerized deployment.

### 1. Clone the Repository

```bash
git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2025-workspace/sebasama/assignment-2.git
```
### 2. Install Dependencies
```bash
go mod tidy
```
### 3. Configure Environment
- Place Firestore service account JSON file at `config/service-account.json`.
- Optionally, set the `PORT` environment variable (default 8080).

## Run the Application
#### Using Go:
```bash
go run main.go
```
#### Using Docker:
```bash
docker compose build
```
-   #### Attached:
    ```bash
    docker compose up assignment-2
    ```
-   #### Detached:
    ```bash
    docker compose up assignment-2 -d
    ```
-   #### Stop service:
    ```bash
    docker compose down assignment-2 
    ```

## API-endpoints
```
/dashboard/v1/registrations/
/dashboard/v1/dashboards/
/dashboard/v1/notifications/
/dashboard/v1/status/
```
### Endpoint 'Registrations'

#### - Request (POST)
```
Method: POST
Path: /dashboard/v1/registrations/
Content type: application/json
```
- **Description:**  
  - Registers a new dashboard configuration indicating which country details and features should be displayed on the dashboard.


- **Example Request Body:**
  ```json
  {
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": true,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "USD", "SEK"]
    }
  }
  
- **Response:**
    ```json
  {
  "id": "v9KIhCCocXgSPwLg8UWN",
  "lastChange": "2025-04-08 12:37:09.4366611 +0200 CEST"
  }
#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/registrations/{id}
Content type: application/json
```
- **Description:**  
  - Retrieves the complete dashboard configuration corresponding to the provided ID.


- **Example Request:**
  - `/dashboard/v1/registrations/v9KIhCCocXgSPwLg8UWN/`


- **Response:**
    ```json 
  {
    "id": "v9KIhCCocXgSPwLg8UWN",
    "country": "Norway",
    "isoCode": "NO",
    "features": {
        "temperature": true,
        "precipitation": true,
        "capital": true,
        "coordinates": true,
        "population": true,
        "area": false,
        "targetCurrencies": [
            "EUR",
            "USD",
            "SEK"
        ]
    },
    "lastChange": "2025-04-08 12:37:09.4366611 +0200 CEST"
  }

#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/registrations/
```
- **Description:**  
  - Returns an array of all dashboard configurations.


- **Request:**
  - `/dashboard/v1/registrations/`


- **Response:**
    ```json 
  [
    {
        "id": "Dq8PjJfuPVQT8YEq8Ms4",
        "country": "Norway",
        "isoCode": "NO",
        "features": {
            "temperature": true,
            "precipitation": true,
            "capital": true,
            "coordinates": true,
            "population": true,
            "area": false,
            "targetCurrencies": [
                "EUR",
                "USD",
                "SEK"
            ]
        },
        "lastChange": "2025-04-08 12:30:46.9248986 +0200 CEST"
    },
    {
        "id": "Gzt5covH1vt8y3QVtnan",
        "country": "NO",
        "isoCode": "",
        "features": {
            "temperature": false,
            "precipitation": false,
            "capital": false,
            "coordinates": false,
            "population": false,
            "area": false,
            "targetCurrencies": null
        },
        "lastChange": "2025-04-08 12:29:57.1243472 +0200 CEST"
    },
    ...



#### - Request (PUT)
```
Method: PUT
Path: /dashboard/v1/registrations/{id}
```
- **Description:**  
  - Replaces the entire dashboard configuration identified by the provided ID and updates lastChange timestamp. 


- **Example Request Body:**
    ```json
    {
    "country": "Norway",
    "isoCode": "NO",
    "features": {
      "temperature": false,
      "precipitation": true,
      "capital": true,
      "coordinates": true,
      "population": true,
      "area": false,
      "targetCurrencies": ["EUR", "SEK"]
      }
  }



- **Response:**
  - Returns 204 No Content
  - Body: empty

#### - Request (PATCH)
```
Method: PATCH
Path: /dashboard/v1/registrations/{id}
Content type: application/json
```
- **Description:**
  - Partially updates the dashboard configuration, modifying only the provided fields and automatically updating the lastChange timestamp.


- **Example Request Body:**
    ```json
    {
    "features": {
      "temperature": false
        }
  }

- **Response:**
    - Returns 204 No Content
    - Body: empty

#### - Request (HEAD)
**Note:** For HEAD requests, the ID parameter is optional. When no ID is provided, the request applies to the entire collection.

```
Method: HEAD
Path: /dashboard/v1/registrations/{id}
```
- **Description:**
  -   Retrieves only the headers for the dashboard configuration with the specified ID. This can be used to verify the existence of the resource and inspect its metadata, without returning body.


- **Example Request:**
  - `/dashboard/v1/registrations/v9KIhCCocXgSPwLg8UWN/`


- **Response:**
  - Returns 204 No Content

#### - Request (DELETE)
```
Method: DELETE
Path: /dashboard/v1/registrations/{id}
```
- **Description:**
  - Deletes the dashboard configuration identified by the provided ID.


- **Example Request:**
    - `/dashboard/v1/registrations/v9KIhCCocXgSPwLg8UWN/`


- **Response:**
    - Returns 204 No Content
    - Body: empty


### Endpoint 'Dashboards'


#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/dashboards/{id}
```
- **Description:**
  - Retrieves a populated dashboard identified by the provided ID.
  

- **Example Request:**
    - `/dashboard/v1/dashboards/v9KIhCCocXgSPwLg8UWN/`


- **Response:**
    ```json
  {
    "country": "Norway",
    "features": {
        "area": 323802,
        "capital": [
            "Oslo"
        ],
        "coordinates": {
            "latitude": 62,
            "longitude": 10
        },
        "population": 5379475,
        "precipitation": 10.428571428571429,
        "targetCurrencies": [
            {
                "base_code": "NOK",
                "time_last_update_utc": "Mon, 07 Apr 2025 00:02:31 +0000",
                "time_next_update_utc": "Tue, 08 Apr 2025 00:23:41 +0000",
                "rates": [
                    {
                        "code": "EUR",
                        "rate": 0.084486
                    },
                    {
                        "code": "USD",
                        "rate": 0.092858
                    },
                    {
                        "code": "SEK",
                        "rate": 0.929191
                    }
                ]
            }
        ],
        "temperature": 2.457142857142857
    },
    "isoCode": "NO",
    "lastRetrieval": "2025-04-08T13:33:40+02:00"
}

### Endpoint 'Notifications'


#### - Request (POST)
```
Method: POST
Path: /dashboard/v1/notifications/
Content type: application/json
```
- **Description:**
  - Registers a new webhook to be invoked when a specific event occurs.


- **Example Request Body:**
    ```json
  {
  "url": "https://localhost:8080/client/",
  "country": "NO",
  "event": "INVOKE"
  }

- **Response:**
  - Status code: 201 Created
      ```json
    {
      "httpCat": "https://http.cat/201",
      "id": "d7i0baIRFVRS6RB5vCJZ"
    }

#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/notifications/{id}
```
- **Description:**
  - Retrieves the details of a specific webhook registration.


- **Example Request:**
  - `/dashboard/v1/notifications/NDDWIy5TM6kaaAYvIN6E/`


- **Response:**
    ```json
  {
    "id": "NDDWIy5TM6kaaAYvIN6E",
    "url": "https://localhost:8080/client/",
    "country": "NO",
    "event": "INVOKE"
  }

#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/notifications/
```
- **Description:**
  - Retrieves an array of all registered webhooks.


- **Request:**
    - `/dashboard/v1/notifications/`


- **Response:**
    ```json
  [
      {
        "id": "0MwxWRziYlRIANVbmcdh",
        "url": "https://example.com/webhook",
        "country": "NO",
        "event": "REGISTER"
      },
      {
        "id": "7DjSRD7FXUjRXVATFGtL",
        "url": "test.com",
        "country": "NO",
        "event": "REGISTER"
      },
      ...
  ]
#### - Request (PATCH)
```
Method: PATCH
Path: /dashboard/v1/notifications/{id}
Content type: application/json
```
- **Description:**
  - Partially updates webhook registration identified by provided ID, and automatically sets a new lastChange timestamp.


- **Example Request Body:**
    ```json
  {
  "country": "NO",
  "url": "https://updated-example.com/webhook"
  }

- **Response:**
  - Returns 204 No Content
  - Body: empty

#### - Request (HEAD)
**Note:** For HEAD requests, the ID parameter is optional. When no ID is provided, the request applies to the entire collection.

```
Method: HEAD
Path: /dashboard/v1/notifications/{id}
```
- **Description:**
  -   Retrieves only the headers for the dashboard configuration with the specified ID. This can be used to verify the existence of the resource and inspect its metadata, without returning body.


- **Example Request:**
  - `/dashboard/v1/notifications/NDDWIy5TM6kaaAYvIN6E/`


- **Response:**
  - Returns 204 No Content

#### - Request (DELETE)
```
Method: DELETE
Path: /dashboard/v1/notifications/{id}
```
- **Description:**
  - Deletes the webhook registration identified by the provided ID.

  
- **Example Request:**
  - `/dashboard/v1/notifications/NDDWIy5TM6kaaAYvIN6E/`


- **Response:**
  - Returns 204 No Content
  - Body: empty

### Endpoint 'Status'


#### - Request (GET)
```
Method: GET
Path: /dashboard/v1/status/
```
- **Description:**
  - Returns the overall system status, including the HTTP status codes for external APIs (REST Countries, Open Meteo, Currency API), the number of registered webhooks, the service version, and the uptime.


- **Request:**
  - `/dashboard/v1/status/`


- **Response:** 
  - Content type: application/json
  - Status code: 200 if OK

      ```json
    {
      "countriesAPI": 200,
      "currencyAPI": 200,
      "openmeteoAPI": 200,
      "notificationresponse": 200,
      "dashboardresponse": 200,
      "webhookssum": 13,
      "version": "v1",
      "uptime": "0d:02h:33m:38s"
    }

## Caching and Purging
## Testing
## Deployment