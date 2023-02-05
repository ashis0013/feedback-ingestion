# feedback ingestion system

## Key Points
- Heterogeneous sources integration: To support heterogeneous sources, we can use a data ingestion layer to handle multiple input formats and transform them into a standard internal structure.
- Push and Pull Integration Model: To handle both push and pull integration models, we can use an API endpoint for push integration and scheduled data retrievers for pull integrations. There will be unique retrievers for each source.
- Metadata Ingestion: To handle metadata ingestion, for each source we will specify a structure of metadata and the actual data will be appended with the ingested feedback.
- Multi-tenancy: To support multi-tenancy, we can use a tenant id in each record to distinguish between different tenants' data.
- Transformation to a uniform internal structure: To ensure data consistency, we will define a common structure for internal representation and our data ingestion layer will perform the transformation task accordingly.

## How to run
- Make sure to install postgres and go runtime 
- Add env variables for username (PG_USER), password (PASS) and database name (DB_NAME)
- Now you can simply run `go run .` at the project root and it will start

## Design
<img width="560px" src="https://user-images.githubusercontent.com/31564734/216831898-52a24f01-31e2-47a9-ab01-975c7316bdee.png"/>

for more details read <a href="https://github.com/ashis0013/feedback-ingestion/files/10611556/Design.Document.1.pdf">this</a>
