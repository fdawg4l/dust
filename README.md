# Dust -PMS7003 + Arduino(ESP8266) Air Quality Sensor

Implements a PMS7003 based air quality sensor.  There are lots of these solutions on the usual sites,
but I couldn't tell if they ran the PMS7003 continuously or cycled the sensor on a clock.  The LED
has a 1 year lifetime, so it'd be a bummer to buy a thing that you'd have to chuck every year.

This sensor implements a RESTful server.  Put your SSID/password in the right place, build, upload, and
the PMS7003 should go into standby (fan off).  Hit the IP with a `GET` on `:80`, and some artisanal,
hand rolled, JSON should be returned.

```
$ curl http://dust.local | jq '.'

[
  {
    "SP_1_0": 11,
    "SP_2_5": 13,
    "SP_10_0": 15,
    "AE_1_0": 15,
    "AE_2_5": 17,
    "AE_10_0": 18
  },
  {
    "SP_1_0": 12,
    "SP_2_5": 13,
    "SP_10_0": 16,
    "AE_1_0": 15,
    "AE_2_5": 17,
    "AE_10_0": 20
  },
  {
    "SP_1_0": 9,
    "SP_2_5": 12,
    "SP_10_0": 15,
    "AE_1_0": 12,
    "AE_2_5": 15,
    "AE_10_0": 18
  },
  {
    "SP_1_0": 10,
    "SP_2_5": 12,
    "SP_10_0": 15,
    "AE_1_0": 13,
    "AE_2_5": 15,
    "AE_10_0": 18
  },
  {
    "SP_1_0": 11,
    "SP_2_5": 13,
    "SP_10_0": 18,
    "AE_1_0": 15,
    "AE_2_5": 17,
    "AE_10_0": 22
  }
]
```

According the documentation, `SP` is _Standard Particles, CF=1_, and `AE` is _Atmospheric Environment_.  `AE_2_5` seems to be the
most useful for tracking the level of hazardous particulate matter in the air.

There are other projects that use the very same hardware.  I just couldn't get them to work or found them difficult to read.
