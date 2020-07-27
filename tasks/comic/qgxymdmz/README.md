## ABOUT

爬虫
http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga.html

## 方式
```text

1. 获取分类列表，分类返回中有总的条目数以及每一本书的ID和章节数
>>>
curl 'http://content.mobgkt.com/api/comic/genres?showInPortal=true' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
  -H 'Origin: http://wap.qgxymdmz.com' \
  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-genres.html?scene=genres' \
  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
  --compressed \
  --insecure
>>>
{
  "code": 0,
  "message": null,
  "data": [
    {
      "genre": "Action",
      "count": 202,
      "sortIndex": 1,
      "showInPortal": true
    }
  ]
}


>>>
curl 'http://content.mobgkt.com/api/comic/list?pageNo=1&pageSize=200&online=true&sortType=1&genre=Action&level=5' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
  -H 'Origin: http://wap.qgxymdmz.com' \
  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-genres.html?scene=genres' \
  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
  --compressed \
  --insecure
<<<
{
  "code": 0,
  "message": null,
  "data": {
    "pageNo": 1,
    "pageSize": 20,
    "pageCount": 2,
    "totalCount": 31,
    "modelList": [
      {
        "identification": "4daca2a0-2355-4a4a-b743-fbcf97217f91",
        "comicTitle": "The God Of High School",
        "alternative": "갓 오브 하이스쿨 ; The God of Highschool ; God of Highschool ; GoH ; 高校之神",
        "genre": "Shounen ; Action ; Adventure ; Comedy ; Martial arts ; Supernatural ; Webtoon",
        "author": "Park yong-je",
        "releaseTime": "2011",
        "comicStatus": "Ongoing",
        "illustration": "http://justdownit.s3.amazonaws.com/ContentManage/comic/the_god_of_high_school/the_god_of_high_school.png",
        "introduction": "From: Easy Going Scans While an island half-disappearing from the face of the earth, a mysterious organization is sending out invitations for a tournament to every skilled fighter in the world. “If you win you can have ANYTHING you want” They’re recruiting only the best to fight the best and claim the title of The God of high school!",
        "mainColor": "74,65,73",
        "chaptersCount": 385,
        "chapters": null,
        "level": 2,
        "online": true
      }
    ]
  }
}



```

```text

curl 'http://content.mobgkt.com/api/comic/find?comicId=4daca2a0-2355-4a4a-b743-fbcf97217f91' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
  -H 'Origin: http://wap.qgxymdmz.com' \
  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-detail.html?Id=4daca2a0-2355-4a4a-b743-fbcf97217f91&scene=genre' \
  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
  --compressed \
  --insecure

{
  "code": 0,
  "message": null,
  "data": {
    "identification": "4daca2a0-2355-4a4a-b743-fbcf97217f91",
    "comicTitle": "The God Of High School",
    "alternative": "갓 오브 하이스쿨 ; The God of Highschool ; God of Highschool ; GoH ; 高校之神",
    "genre": "Shounen ; Action ; Adventure ; Comedy ; Martial arts ; Supernatural ; Webtoon",
    "author": "Park yong-je",
    "releaseTime": "2011",
    "comicStatus": "Ongoing",
    "illustration": "http://justdownit.s3.amazonaws.com/ContentManage/comic/the_god_of_high_school/the_god_of_high_school.png",
    "introduction": "From: Easy Going Scans While an island half-disappearing from the face of the earth, a mysterious organization is sending out invitations for a tournament to every skilled fighter in the world. “If you win you can have ANYTHING you want” They’re recruiting only the best to fight the best and claim the title of The God of high school!",
    "mainColor": "74,65,73",
    "chaptersCount": 385,
    "chapters": [
      {
        "comicChapterTitle": "ch.1",
        "chapterIndex": 1,
        "chapterPartIndex": null,
        "comicId": "4daca2a0-2355-4a4a-b743-fbcf97217f91",
        "imgLink": null,
        "lastChapter": false,
        "showIndex": 1
      }
    ],
    "level": 2,
    "online": true
  }
}

```

```text
按照经验看，chapterIndex 需要一直循环下去，不一定和摘要中给的页数是一样的

curl 'http://content.mobgkt.com/api/comic/chapter/list?comicId=4daca2a0-2355-4a4a-b743-fbcf97217f91&&pageNo=1&pageSize=100&chapterIndex=1' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
  -H 'Origin: http://wap.qgxymdmz.com' \
  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-inside.html?Id=4daca2a0-2355-4a4a-b743-fbcf97217f91&chapter=1&scene=genre' \
  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
  --compressed \
  --insecure

curl 'http://content.mobgkt.com/api/comic/chapter/list?comicId=4daca2a0-2355-4a4a-b743-fbcf97217f91&&pageNo=1&pageSize=100&chapterIndex=206' \
  -H 'Connection: keep-alive' \
  -H 'Accept: application/json, text/javascript, */*; q=0.01' \
  -H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Mobile Safari/537.36' \
  -H 'Origin: http://wap.qgxymdmz.com' \
  -H 'Referer: http://wap.qgxymdmz.com/pinkmanga/l/mobile/manga-inside.html?Id=4daca2a0-2355-4a4a-b743-fbcf97217f91&chapter=2&scene=genre' \
  -H 'Accept-Language: zh-CN,zh;q=0.9,en-SG;q=0.8,en;q=0.7' \
  --compressed \
  --insecure

{
  "code": 0,
  "message": null,
  "data": {
    "pageNo": 1,
    "pageSize": 100,
    "pageCount": 1,
    "totalCount": 55,
    "modelList": [
      {
        "comicChapterTitle": "ch.1",
        "chapterIndex": 1,
        "chapterPartIndex": 1,
        "comicId": "4daca2a0-2355-4a4a-b743-fbcf97217f91",
        "imgLink": "http://justdownit.s3.amazonaws.com/ContentManage/comic/the_god_of_high_school/the_god_of_high_school.1.1.jpg",
        "lastChapter": false,
        "showIndex": null
      }
    ]
  }
}
```