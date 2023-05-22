# google-play-parser

Clone of the [library](https://www.npmjs.com/package/google-play-scraper)

- [Quick example](#quick-example)
- [Features](#features)
  - [App](#app)
  - [Similar](#similar)
  - [Developer](#developer)
  - [Search](#search)
  - [Data Safety](#data-safety)
  - [Permissions](#permissions)
  - [Suggest](#suggest)
  - [Reviews](#reviews)

## Quick Example

```go
package main

import (
    "context"
    "log"

    gpp "github.com/bots-house/google-play-parser"
)

func main() {
    collector := gpp.New()

    app, err := gpp.App(context.Background(), gpp.ApplicationSpec{
        AppID: "com.tinder",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Println(app)
}
```

## Features

## App

Method for parsing app data from google store

**Parameters**

- `app-id` - platform readable id such **com.tinder** [required]
- `lang`
- `country` in ISO format

**Example result**

```json
{
  "app_id": "com.tinder",
  "url": "https://play.google.com/store/apps/details?gl=us&hl=en&id=com.tinder",
  "title": "Tinder Dating app. Meet People",
  "description": "~~app description~~",
  "summary": "Dating your way! Match, chat, and make new friends for dates or find friends",
  "installs": "100,000,000+",
  "min_installs": 100000000,
  "max_installs": 361942145,
  "currency": "USD",
  "price_text": "Free",
  "free": true,
  "score": 3.6711967,
  "score_text": "3.7",
  "ratings": 5799246,
  "reviews": 250755,
  "histogram": {
    "1": 1204726,
    "2": 314420,
    "3": 607810,
    "4": 728241,
    "5": 2944028
  },
  "available": true,
  "offers_iap": true,
  "iap_range": "$0.99 - $299.99 per item",
  "android_version": "7.0",
  "android_version_text": "7.0",
  "developer": "8070166968320699506",
  "developer_id": "8070166968320699506",
  "developer_email": "help@gotinder.com",
  "developer_website": "https://tinder.com",
  "developer_address": "Tinder\n8833 W. Sunset Blvd.\nWest Hollywood, CA 90069",
  "privacy_policy": "https://policies.tinder.com/privacy",
  "genre": "Dating",
  "genre_id": "DATING",
  "icon": "https://play-lh.googleusercontent.com/fDpoqIbZ884ylRnMK8Lx9Fu4DsLQk5yt4f9WkxeOAPpGnzc9BTi_YKkMsLvoMdx7Uzg",
  "header_image": "https://play-lh.googleusercontent.com/fDpoqIbZ884ylRnMK8Lx9Fu4DsLQk5yt4f9WkxeOAPpGnzc9BTi_YKkMsLvoMdx7Uzg",
  "screenshots": [
    "https://play-lh.googleusercontent.com/YjX6U0xrpDX6p9bRqfyaiIcr8LmWJQjKpjEhofh54p3T9MZq8y-bHBpZTUDKDqrh",
    "https://play-lh.googleusercontent.com/WWJE1wosHL4uo1qX6KAmOAP3N_V4RCyK6bMJO1KaKSWc3hcKWm8INy0KO4PORnSnnBc",
    "https://play-lh.googleusercontent.com/Anwn4H8ay1LJFx-uDoVqCDLeBydcK2THS0OeH44FRV0I4H7Zi1adLwqF3TLckK94knP_",
    "https://play-lh.googleusercontent.com/CKuVZ-0vtkTf3wWG6_l8LHlN8Ee4thkjIHahZ-UAxy97B4UoekWrlY4TxcQXYauVqTI",
    "https://play-lh.googleusercontent.com/vSCIDKLJgTmP_Sww65mA7cmIPU89oJQe4Ufy6Toiaayq7i1hoxR8YgL5ctnq1HLJtGg",
    "https://play-lh.googleusercontent.com/aT9_hJ8IXbbMY-Hjbp6qFZSLEsh-gleyT0L1pJMHlXpCq-f-JkHechjM2BBTVA6GFyzS",
    "https://play-lh.googleusercontent.com/b3MfPeeCBKisHMmImXD6LDRPtr7hly342AI6wik91NGEFpQBzZvCQePmbljOJxncjw",
    "https://play-lh.googleusercontent.com/EhuGna9qCDVYvGykjR0BV6rkESFKDAu6zYxqCp2rMAlWmesbYUpMyjD-8rU68yQh1A"
  ],
  "content_rating": "Mature 17+",
  "ad_supported": true,
  "released": "Jul 15, 2013",
  "updated": 1684188049000,
  "version": "14.9.0",
  "recent_changes": "Bug fixes and improvements"
}
```

---

## Similar

Method for parsing app data which similar for requested

**Parameters**

- `app-id` - platform readable id such **com.tinder** [required]
- `lang`
- `country` in ISO format
- `count`
- `full` if false parse only common data

**Returns array of apps**

---

## List

Method which parse list of apps

**Parameters**

- `age` - possible [values](#list-age)
- `category` - possible [values](#list-category)
- `collection` - possible [values](#list-collection)
- `lang`
- `country` in ISO format
- `count`
- `full` if false parse only common data

**Returns array of apps**

---

## Developer

Method which parse developer apps data

**Parameters**

- `dev-id` - developer id numeric or full like
- `count`
- `lang`
- `country` in ISO format
- `full` if false parse only common data

**Returns array of apps**

---

## Search

Method for parsing apps data by some query

**Parameters**

- `query` - search params
- `count`
- `price` - possible [values](#search-price)
- `lang`
- `country` in ISO format
- `full` if false parse only common data

## **Returns array of apps**

## Data Safety

Method which parse app data safety

**Parameters**

- `app-id` - platform readable id such **com.tinder** [required]
- `lang`

**Example result**

```json
{
  "shared_data": [
    {
      "data": "User IDs",
      "optional": false,
      "purpose": "Advertising or marketing, Account management",
      "type": ""
    },
    {
      "data": "Installed apps",
      "optional": false,
      "purpose": "Advertising or marketing",
      "type": ""
    },
    {
      "data": "Crash logs",
      "optional": false,
      "purpose": "Analytics",
      "type": ""
    }
  ],
  "collected_data": [
    {
      "data": "Name",
      "optional": true,
      "purpose": "App functionality, Developer communications, Advertising or marketing",
      "type": ""
    },
    {
      "data": "Purchase history",
      "optional": true,
      "purpose": "Account management",
      "type": ""
    },
    {
      "data": "Other in-app messages",
      "optional": false,
      "purpose": "Developer communications, Fraud prevention, security, and compliance",
      "type": ""
    },
    {
      "data": "Contacts",
      "optional": true,
      "purpose": "App functionality",
      "type": ""
    },
    {
      "data": "Other actions",
      "optional": false,
      "purpose": "App functionality, Analytics, Fraud prevention, security, and compliance",
      "type": ""
    },
    {
      "data": "Crash logs",
      "optional": true,
      "purpose": "App functionality, Analytics",
      "type": ""
    },
    {
      "data": "Device or other IDs",
      "optional": false,
      "purpose": "App functionality, Analytics, Advertising or marketing, Fraud prevention, security, and compliance, Personalization, Account management",
      "type": ""
    }
  ],
  "privacy_policy_url": "http://www.jamcity.com/privacy",
  "security_practice": [
    {
      "description": "Your data isn’t transferred over a secure connection",
      "practice": "Data isn’t encrypted"
    },
    {
      "description": "The developer provides a way for you to request that your data be deleted",
      "practice": "You can request that data be deleted"
    }
  ]
}
```

---

## Permissions

Method which parse app permissions data

**Parameters**

- `app-id` - platform readable id such **com.tinder** [required]
- `lang`
- `full` if true parse description of permission

**Example result**

```json
[
  {
    "type": "Phone"
  },
  {
    "type": "Wi-Fi connection information"
  },
  {
    "type": "Device & app history"
  },
  {
    "type": "Device ID & call information"
  },
  {
    "type": "Storage"
  },
  {
    "type": "Photos/Media/Files"
  }
]
```

---

## Suggest

Method for parsing suggest by search query

**Parameters**

- `query`- search query
- `lang`
- `country` in ISO format

**Example result**

```json
["paypal", "paycom", "paylocity", "paychex", "payrange"]
```

---

## Reviews

Method which parse app reviews

**Parameters**

- `app-id` - platform readable id such **com.tinder** [required]
- `lang`
- `country` in ISO format
- `count`
- `sort` - [1, 2, 3]

**Example result**

```json
[
  {
    "id": "1ed5a1d7-ce39-4179-8651-ac8b05084eb8",
    "url": "https://play.google.com/store/apps/details?id=com.sgn.pandapop.gp&reviewId=1ed5a1d7-ce39-4179-8651-ac8b05084eb8",
    "summary": "Having the same issues with not being able to pass the upper levels without purchase. Currently on level 6268. No way to complete the level. The level offers no boosters and even seems to change the ball colors to insure a loss. If the game developers are trying to slow players down on higher levels, it's working. Seems like once you start purchasing boosters to pass a level, the games want to perpetuate in app sales. This player is also losing interest in playing the game altogether. Too bad.",
    "score": 3,
    "score_text": "3.00",
    "user_image": "https://play-lh.googleusercontent.com/a/AGNmyxZY6n0JZceC0vwo3_ErqR_IFjaVVJxpy3Q1oc_m=mo",
    "user_name": "Sue T",
    "version": "12.3.103",
    "date": "2023-04-21T17:16:35.33+03:00",
    "reply_text": "Hi Sue, we completely understand how you feel about this, we're working hard to improve our game and your satisfaction is our main priority, rest assured that we'll send your comments about the hard levels over to our team for further consideration. Thanks for your feedback!",
    "reply_date": "2023-04-26T19:06:44.869+03:00",
    "criteria": [
      {
        "criteria": "vaf_games_genre_claw",
        "rating": 2
      }
    ],
    "tumbs_up": 0
  },
  ...
]
```

---

---

## List age

- "AGE_RANGE1"
- "AGE_RANGE2"
- "AGE_RANGE3"

## List collection

- "TOP_FREE"
- "TOP_PAID"
- "GROSSING"

## List category

- "APPLICATION"
- "ANDROID_WEAR"
- "ART_AND_DESIGN"
- "AUTO_AND_VEHICLES"
- "BEAUTY"
- "BOOKS_AND_REFERENCE"
- "BUSINESS"
- "COMICS"
- "COMMUNICATION"
- "DATING"
- "EDUCATION"
- "ENTERTAINMENT"
- "EVENTS"
- "FINANCE"
- "FOOD_AND_DRINK"
- "HEALTH_AND_FITNESS"
- "HOUSE_AND_HOME"
- "LIBRARIES_AND_DEMO"
- "LIFESTYLE"
- "MAPS_AND_NAVIGATION"
- "MEDICAL"
- "MUSIC_AND_AUDIO"
- "NEWS_AND_MAGAZINES"
- "PARENTING"
- "PERSONALIZATION"
- "PHOTOGRAPHY"
- "PRODUCTIVITY"
- "SHOPPING"
- "SOCIAL"
- "SPORTS"
- "TOOLS"
- "TRAVEL_AND_LOCAL"
- "VIDEO_PLAYERS"
- "WATCH_FACE"
- "WEATHER"
- "GAME"
- "GAME_ACTION"
- "GAME_ADVENTURE"
- "GAME_ARCADE"
- "GAME_BOARD"
- "GAME_CARD"
- "GAME_CASINO"
- "GAME_CASUAL"
- "GAME_EDUCATIONAL"
- "GAME_MUSIC"
- "GAME_PUZZLE"
- "GAME_RACING"
- "GAME_ROLE_PLAYING"
- "GAME_SIMULATION"
- "GAME_SPORTS"
- "GAME_STRATEGY"
- "GAME_TRIVIA"
- "GAME_WORD"
- "FAMILY"

## Search price

- all
- free
- paid
