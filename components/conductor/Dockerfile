#
# Copyright (C) 2018 Nalej Group - All Rights Reserved
#

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

# Set nalej working directory
# RUN mkdir /nalej
# WORKDIR /nalej

# Now just add the binary
COPY conductor /nalej/

ENTRYPOINT ["./nalej/conductor"]