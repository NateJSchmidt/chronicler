# Event Sourcing

Event sourcing is too large of a topic to cover here.  There are numerous resources freely available online.  If you are looking for a place to start, I'd recommend watching some of Greg Young's talks ([this](https://www.youtube.com/watch?v=JHGkaShoyNs) is a good place to start).

Chronicler is an opinionated event log.  As an event log, its functional purpose is to store events, allow services to retrieve events, and provide a mechanism for services to be notified that events have happened.  In order to accomplish this goal, Chronicler is built on the following event sourcing principles:

1. Event stream names are of the form `<noun>:<verb>-<id>` (e.g. `user:login-7a99bfb8-6c8f-47b5-ba84-8bca6f073c97`).
   1. Everything to the left of the first hyphen is called the _category_.  A _category_ represent a unique type of event that is provided by a single service.
   1. The unique identifier is everything to the right of the first hyphen (e.g. `7a99bfb8-6c8f-47b5-ba84-8bca6f073c97` in our example).  The unique identifier can be used to correlate related events across multiple _categories_.
1. All unique identifiers are UUIDs that adhere to the [UUID4](https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)) standard.
1. Each event in the system belongs to an event stream and can be identified by a unique identifier (which is different than the unique identifier of the stream name).
