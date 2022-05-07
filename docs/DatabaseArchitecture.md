# Database architecture

The database architecture of Chronicler is based on the following principles:

1. Scales over time as the number of events and event streams grow.
1. Adheres to microservice best practices such that only a single service writes to a specific table.  In other words, there will not be multiple, different services all reading and writing to the same database table.
1. Allows for service horizontal scalability such that a single service **can** be horizontally scaled and all of the replicas will be able to read and write to the same table.

## Database tables

To accomplish these goals, Chronicler creates a separate database table for each unique category of event streams.  Each table is named after the category.  For example, the `user:login` category of streams would be stored in the `user_login` table.

If Chronicler receives a request to record an event on a stream that belongs to a category that does not exist in the database, it will create a new table for that category of streams.  This naturally supports the growth and addition of new services into the event sourcing system over time without having to modify the database schema.

## Database schema

Each table contains the following columns:
| Column Name | Type   | Description |
| ----------- | ------ | ----------- |
| id          | string | This is the unqiue identifier of the event.  It is a UUID4. |
| time        | datetime | This is the time that the event was recorded in Chronicler. |
| event       | jsonb    | This is the event payload. |
| streamName  | string | This is the name of the event stream. |
| position    | integer | This is the position of the event in the event stream. It is also the unique identifier for the table. |
| version     | string | This is the version of the event payload.  It is likely a semver based string. |
