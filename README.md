# docker-soju
Docker image for the soju irc bouncer, with a custom helper app to configure DB using environment variables

## Example usage

    docker run shanedabes/soju:latest -p 6667:6667
      -e SOJU_USER_1_NAME=shane \
      -e SOJU_USER_1_PASSWORD=mypass \
      -e SOJU_USER_1_NETWORK_1_NAME=freenode \
      -e SOJU_USER_1_NETWORK_1_SERVER=chat.freenode.net \
      -e SOJU_USER_1_NETWORK_1_NICK=shane \
      -e SOJU_USER_1_NETWORK_1_PASSWORD=nickpass \
      -e SOJU_USER_1_NETWORK_1_CHANNELS=#soju,#go-nuts \
      -e SOJU_USER_1_NETWORK_2_NAME=synirc # ...etc \
      -e SOJU_USER_2_NAME=otheruser # ...etc \
      -e SOJU_TRUST_ADD=delta.uk.eu.synirc.net:6697 # for problematic servers that need to be added to trust
