# triple-s

This project is designed to implement a simplified version of S3 (Simple Storage Service) object storage. This tool will provide a REST API that allows clients to interact with the storage system, offering core functionalities such as creating and managing buckets, uploading, retrieving, and deleting files, as well as handling object metadata.
## Flags
#### port
Port to work on.
#### dir
Directory to store buckets and objects.
#### help
Returns help message
## Handlers
### buckets
#### 1. Create a Bucket:
HTTP Method: PUT

Validate the bucket name to ensure it meets Amazon S3 naming requirements, ensure the bucket name is unique across the entire storage system, create a new entry in the bucket metadata storage, and return a 200 OK status code and details of the created bucket, or an appropriate error.
#### 2. List All Buckets:
HTTP Method: GET

Read the bucket metadata from the storage (e.g., a CSV file), return an XML response containing a list of all matching buckets, including metadata like creation time, last modified time, etc, and respond with a 200 OK status code and the XML list of buckets.
#### 3. Delete a Bucket:
HTTP Method: DELETE

Check if the specified bucket exists by looking it up in the bucket metadata storage, ensure the bucket is empty (no objects are stored in it) before deletion. If the bucket exists and is empty, remove it from the metadata storage, and return a 204 No Content status code if the deletion is successful, or an error message in XML format if the bucket does not exist or is not empty.
### objects
#### 1. Upload a New Object:

HTTP Method: PUT

Verify if the specified bucket exists by reading from the bucket metadata, validate the object key, save the object content to a file in a directory named after the bucket, store object metadata in a CSV file, and respond with a 200 status code or an appropriate error message if the upload fails.

#### 2. Retrieve an Object:

HTTP Method: GET

Verify if the bucket exists, check if the object exists and return the object data or an error.

#### 3. Delete an Object:

HTTP Method: DELETE

Verify if the bucket and object exist, delete the object and update metadata, and respond with a 204 No Content status code or an appropriate error message.