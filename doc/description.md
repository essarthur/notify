# Description about service

Usage
https://github.com/appleboy/gorush#run-grpc-service


**ESS-BRIDGE-NOTIFICATIONS**


**Purpose:**
1. prepare push notifications for Essentia iOS,Android and Desktop applications;
2. send push notifications to [APNS](https://developer.apple.com/notifications/) and [FCM](https://firebase.google.com/docs/cloud-messaging/);
3. provide functionality for switching user's account via REST API endpoint, applying user's configs;
4. store user's configs (id, address, push notifications type, time of sending, etc);
5. logging. It is able to log information (info, warn, debug, err, etc);
6. analytics (creating reports, charts, analytics, etc).

Bridge stores SSL-Certificate for signing messages before sending to the APNS.

**PostrgeSQL** DB is used for storing data. DB backup service is provided.
**Redis** DB is used for message queue, cache.

Third-party services: [Gorush](https://github.com/appleboy/gorush).

**General flow:**


![Flow](Flow.png)

**Struct:**

![Struct](Struct.png)


**Remote part:**
1. gorush instancies - proxy between ESS-Bridge-Notification and APNS/FCM.


**Main components:**
1. **Android and iOS notification creators** (NC). Purpose: create specific notifications (pre-notification for gorush) using input data.
2. **Wallet notification manager** (WNM):
	1. iteracts with NC, prepare data for NC, receive notification from NC and transfer it to **Delivery worker switcher** (DWS).
	2. iteracts with **Device token/info handler**, takes device token and account info from DB.
	3. proccess input request from **Request parser** (RP).
	4. send requests to **Bridge-Wallet adapter** for info about concrete wallet status (balance, TX info, etc).
3. **Device token/info handler** (DTH):
	1. iteracts with DB wrapper. Request/procced data from/to DB.
	2. handles input requests from Device registration API, proccess device token/info. Provide info for WNM.
4. **DB Wrapper** - provide simplified access to the DB.

