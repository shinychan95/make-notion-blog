# make-notion-blog

> 막상 열심히(reverse 엔지니어링) 만들었는데, 사용성을 고려하니 정작 나만 쓸 것 같은 도구이다(🫠)

Notion 에서 글을 작성하고, 곧바로 블로그로 원터치로 배포하기 위한 도구이다.

사용 편의성을 가지는 부분은 Notion (복사 붙여넣기, markdown 추출)로 원활하지 않던, (표, 이미지, 기타 컴포넌트) 등을 편리하게 옮길 수 있다.

## 사용법

![prototype](https://user-images.githubusercontent.com/39409255/233854586-3d1705be-4916-4416-ad2f-e8fb8fcdb512.gif)

1. Notion 내 새로운 페이지를 만들고, 해당 페이지 템플릿을 `블로그 포스팅 캘린더` 로 설정해야 한다.
   - 해당 페이지는 `collection_view_page` 타입이다.
   - 블로그 포스팅 메타 데이터에 맞춰 아래 property 들이 연동된다.
      - Status : `Drafting` 상태의 페이지들만이 마크다운 변환된다.
      - Date : _(현재 프로그램 실행 시간으로 동작)_
      - Categories : 블로그 내 Catogories 에 해당하며, 첫번째부터 첫 depth 에 해당한다.
      - Tags : 실제 블로그 내에는 소문자로 입력된다.
2. 이미지 다운로드를 위해 Notion Integration 생성 및 해당 페이지에 Integration 을 연결한다.
   - https://www.notion.so/my-integrations 자신의 workspace 에 associated 하도록 생성한다.
   - Notion 페이지 내에서 ( ... -> 연결 -> Integration 선택) 과정을 수행한다.
3. 웹페이지를 통해 해당 페이지에 접속하고, URL 을 통해 페이지 ID 를 가져온다.
   - https://www.notion.so/shinychan95/806a2a5ddce84729916a387f939bc82b?v=21e14706e4f343bb852e6c3fde98c06a
   - URL 뒷 파라미터 v 는 view type 을 말한다. 무시해도 된다.
4. 프로그램 실행을 위한 설정 파일을 작성한다. (config.json)
    ```json
    {
      "db_path": "/Users/user/Library/Application Support/Notion/notion.db",
      "api_key":"secret_...",
      "post_directory": "/Users/user/github/shinychan95.github.io/_posts",
      "image_directory": "/Users/user/github/shinychan95.github.io/assets/pages",
      "root_id": "806a2a5d-dce8-4729-916a-387f939bc82b"
    }
    ```

## 기술적인 특색
- Notion 어플리케이션이 offline 을 위해 사용하는 cache DB(SQLite) 내 데이터를 파싱하여 마크다운을 생성한다.
- 이미지의 경우, Notion API 를 통해 static image 경로를 받아오고, 이를 다운로드 받아 저장한다.
- Reverse 엔지니어링을 통해 Notion 내 block 데이터들을 type 별로 가져와 마크다운으로 변환한다.

## 제약 사항들
1. Notion 이미지를 위해 API 연동을 위해 Integration 생성 및 연결을 해주어야 한다.
2. Notion 내 블로그 template 의 페이지를 기준으로 만들어졌다.
   - 특히 Status property 가 Drafting 인 페이지들만 마크다운 변환하도록 현재 로직이 동작한다.
3. 블로그의 경우, github pages 로 구동되는 블로그와 크게 의존성을 가진다.
    - 특히 [Chirpy Jekyll Theme](https://github.com/cotes2020/jekyll-theme-chirpy) 테마를 사용하지 않는다면 커스터마이징이 필요할 수 있다.
4. Mac OS 기준으로 작업한 것이므로 Windows 에서 정상 작동하지 않을 가능성이 크다.
