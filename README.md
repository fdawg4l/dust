# Dust:  PMS7003 + Arduino(ESP8266) Air Quality Sensor

Implements a PMS7003 and DHT11 based air quality sensor.  There are lots of
these solutions on the usual sites, but I couldn't tell if they ran the PMS7003
continuously or cycled the sensor on a clock.  The laser LED has a 1 year lifetime,
so it'd be a bummer to buy a thing that you'd have to chuck every year.

This sensor implements a RESTful server.  Put your SSID/password in the right
place, build, upload, and the PMS7003 should go into standby (fan off).  Hit
the IP with a `GET` on `:80`, and some artisanal hand rolled JSON should be
returned.

```
$ curl http://dust.local | jq '.'
{
  "t": {
    "humidity_P": 59,
    "temp_F": 62.6
  },
  "a": [
    {
      "SP_1_0": 7,
      "SP_2_5": 9,
      "SP_10_0": 12,
      "AE_1_0": 10,
      "AE_2_5": 12,
      "AE_10_0": 15
    },
    {
      "SP_1_0": 7,
      "SP_2_5": 9,
      "SP_10_0": 10,
      "AE_1_0": 10,
      "AE_2_5": 12,
      "AE_10_0": 14
    },
    {
      "SP_1_0": 9,
      "SP_2_5": 9,
      "SP_10_0": 10,
      "AE_1_0": 12,
      "AE_2_5": 12,
      "AE_10_0": 14
    },
    {
      "SP_1_0": 9,
      "SP_2_5": 9,
      "SP_10_0": 12,
      "AE_1_0": 12,
      "AE_2_5": 12,
      "AE_10_0": 15
    },
    {
      "SP_1_0": 9,
      "SP_2_5": 9,
      "SP_10_0": 13,
      "AE_1_0": 12,
      "AE_2_5": 12,
      "AE_10_0": 17
    }
  ]
}
```

According the documentation, `SP` is _Standard Particles, CF=1_, and `AE` is _Atmospheric Environment_.  `AE_2_5` seems to be the
most useful for tracking the level of hazardous particulate matter in the air.

## Grafana + influxdb

There's a `docker-compose.yml` and an accompanying `dustd` _microservith_ which periodically pulls data from the sensor and posts it to influxdb.  To build/install it, just do the following
```
$ make && make docker && docker-compose up -d && ./db.sh
```

Obviously you need `docker` and `docker-compose` installed.  Google is your friend if you don't already have that.  Also, you won't get far without `go`.  Again, _le Googs_ will help.

And `db.sh` will probably not help you unless your docker host is called `docker.chaos.local` like mine is.  Basically create the db on your host using curl like `db.sh` is doing.  Also, take a look at `config.json`.  You probably don't need to change anything.  Rather, add whatever changes you want to make as environment variables in `env.dustd`.

In the end you'll get pretty graphs like the following.



## Resources and dependencies
http://aqicn.org/sensor/pms5003-7003/

http://www.handsontec.com/pdf_learn/esp8266-V10.pdf

https://github.com/fu-hsi/PMS

https://github.com/adafruit/DHT-sensor-library


