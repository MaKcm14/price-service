## This is a description of the service eficiency.

# The time of responses.

Due to all markets need to parse the time of the response is clearly depends on the speed of the

browser instance. For example, as you start the service and send the first request you'll get the response

from one market just after the *time = 5 + time_of_parsing*. It's really important to understand that the

requests connected with the next values of the query params *amount=max* and *no-image=0* are slow.

So if you need to have the fast response from the markets use the *amount=min* and *no-image=1*

