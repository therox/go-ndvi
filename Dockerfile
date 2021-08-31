# from alpine

# RUN apk --update gdal gdal-dev pkgconfig g++ go

FROM debian
RUN apt update -y && apt upgrage -y && apt install libgdal-dev gdal-bin

CMD [ "/bin/bash" ]