## Nigerian Universities API Documentation

### Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)

### Overview

This API provides access to a list of Nigerian universities' scraped from the official website of the National Universities Commission (NUC) - https://www.nuc.edu.ng/.

### Prerequisites

- [Go](https://go.dev/doc/install) installed
- [Docker](https://www.docker.com/get-started/) installed (optional)

### Installation

To run the API locally, follow these steps:

1. Clone the repository from Github:

```bash
git clone https://github.com/coolpythoncodes/nigerian-universities.git
```

2. Navigate to the project directory:

```bash
cd nigerian-universities
```

3. Run the application:

```bash
go run main.go
```

The application will be served on port `8080`

### Endpoints

The App provides the following endpoints:

GET / fetches all the Nigerian universities (Federal, State and Private)

Example

```bash
curl http://localhost:8080
```

Response

```json
error: false,
message: "success"
data: [
    {
        name: "Abubakar Tafawa Balewa University, Bauchi",
        vice_chancellor: "Professor M A Abdulazeez",
        year_of_establishment: "1988",
        type: "Federal",
        url: "https://www.atbu.edu.ng"
    },
    {
        name: "Ahmadu Bello University, Zaria",
        vice_chancellor: "Professor Kabir Bala",
        year_of_establishment: "1962",
        type: "Federal",
        url: "https://www.abu.edu.ng"
    },
    // more universities with also State and Private
]
```


GET /federal fetches all the Nigerian Federal universities

Example

```bash
curl http://localhost:8080/federal
```

Response

```json
error: false,
message: "success"
data: [
    {
        name: "Abubakar Tafawa Balewa University, Bauchi",
        vice_chancellor: "Professor M A Abdulazeez",
        year_of_establishment: "1988",
        type: "Federal",
        url: "https://www.atbu.edu.ng"
    },
    {
        name: "Ahmadu Bello University, Zaria",
        vice_chancellor: "Professor Kabir Bala",
        year_of_establishment: "1962",
        type: "Federal",
        url: "https://www.abu.edu.ng"
    },
    // more federal universities 
]
```

GET /state fetches all the Nigerian State universities

Example

```bash
curl http://localhost:8080/state
```

Response

```json
error: false,
message: "success"
data: [
    {
        name: "Abia State University, Uturu",
        vice_chancellor: "Prof. Onyemachi M. Ogbulu",
        year_of_establishment: "1981",
        type: "State",
        url: "https://www.abiastateuniversity.edu.ng"
    },
    {
        name: "Adamawa State University Mubi",
        vice_chancellor: "Prof (Mrs) Kaletapwa Farauta",
        year_of_establishment: "2002",
        type: "State",
        url: "https://www.adsu.edu.ng"
    },
    // more state universities 
]
```


GET /private fetches all the Nigerian Private universities

Example

```bash
curl http://localhost:8080/private
```

Response

```json
error: false,
message: "success"
data: [
    {
        name: "Achievers University, Owo",
        vice_chancellor: "Professor Samuel Aje",
        year_of_establishment: "2007",
        type: "Private",
        url: "https://www.achievers.edu.ng"
    },
    {
        name: "Adeleke University, Ede",
        vice_chancellor: "Prof. Samuel E Alao",
        year_of_establishment: "2011",
        type: "Private",
        url: "https://www.adelekeuniversity.edu.ng"
    },
    // more private universities 
]
```


