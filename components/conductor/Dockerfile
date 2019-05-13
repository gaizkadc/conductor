#
# Copyright (C) 2018 Nalej Group - All Rights Reserved
#

# iron/go is the alpine image with only ca-certificates added
FROM iron/go

# Set nalej working directory
# RUN mkdir /nalej
# WORKDIR /nalej

COPY conductor /nalej/

# Now just add the binary
ENTRYPOINT ["./nalej/conductor"]