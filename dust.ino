
#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <SoftwareSerial.h>
#include "PMS.h"

#define HOSTNAME "dust"

const char* ssid = "ShitakeMushrooms";
const char* password = "XXXXXXXX";
char buffer[512] = "";

ESP8266WebServer server(80);
MDNSResponder mdns;

// Pin silkscreen vs reality
#define D2 4
#define D3 0
#define D4 2

// The fan takes 30s to stabilize and give useful readings.
#define PMS_FAN_DELAY 30000

#define SAMPLES 5

// RX - D2(GPIO4); Tx - D3(GPIO0)
// RX, TX
SoftwareSerial pmsSerial(D2,D3);
PMS pms(pmsSerial);

void setup() {
    // Poke the LED off
    pinMode(D4, OUTPUT);
    digitalWrite(D4, HIGH);

    Serial.begin(115200);

    // They say this thing supports i2c, but the pdf doesn't mention it.
    pmsSerial.begin(9600);

    WiFi.disconnect();
    WiFi.hostname(HOSTNAME);
    WiFi.begin(ssid, password);

    Serial.print("Attempting to connect to WPA SSID: ");
    Serial.println(ssid);

    // Wait for connection
    while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
    }

    // Dump the connection info
   Serial.println("");
   Serial.print("Connected to ");
   Serial.println(ssid);
   Serial.print("IP address: ");
   Serial.println(WiFi.localIP());

    // Start the MDNS
    if (mdns.begin(HOSTNAME)) {
       Serial.println("MDNS responder started");
    }

    // Add service to the MDNS
    mdns.addService("http", "tcp", 80);
    pms.passiveMode();
    pms.sleep();

    server.on("/", []() {
      size_t len = SAMPLES;
      PMS::DATA data[len];

      // Gross hand rolling of json
      if (sample(data, len)) {
        uint8_t idx;
        char *buf = buffer;

        // start with an array
        buf[0] = '[';
        buf++;

        // add an object per item
        for (idx = 0; idx < len; idx++) {

          int n = sprintf(buf,
           "{\"SP_1_0\": %d,"
           "\"SP_2_5\": %d,"
           "\"SP_10_0\": %d,"
           "\"AE_1_0\": %d,"
           "\"AE_2_5\": %d,"
           "\"AE_10_0\": %d}",
           data[idx].PM_SP_UG_1_0,
           data[idx].PM_SP_UG_2_5,
           data[idx].PM_SP_UG_10_0,
           data[idx].PM_AE_UG_1_0,
           data[idx].PM_AE_UG_2_5,
           data[idx].PM_AE_UG_10_0);

          buf = buf + n;

          // close out the array or object
          if (idx + 1 == len) {
              buf[0] = ']';
              buf[1] = 0;
          } else {
              buf[0] = ',';
              buf++;
          }
        }

        server.send(200, "text/plain", buffer);
      } else {
        server.send(503, "text/plain", "oops\n");
      }
    });

    server.begin();
}

void loop() {
  server.handleClient();
}


// Sample will get a sample of data every 5s for the number of len.
// We turn on the light, wake up the sensor, wait 30s for the fan
// to do fan things, then start asking for data.  Then we turn the
// sensor off, and turn off the light.
bool sample(PMS::DATA *out, size_t len) {
  uint8_t idx = 0;
  bool rv = true;
  digitalWrite(D4, LOW);
  pms.wakeUp();
  delay(PMS_FAN_DELAY);

  for(idx = 0; idx < len; idx++) {
    pms.requestRead();
    digitalWrite(D4, LOW);
    if (!pms.readUntil(out[idx])) {
      rv = false;
      goto out;
    }

    digitalWrite(D4, HIGH);
    delay(5000);
  }

out:

  pms.sleep();
  digitalWrite(D4, HIGH);

  return rv;
}
