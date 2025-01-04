import json
from curl_cffi import requests

MM_ORIGIN: str = "https://megamarket.ru"
COLLECTION_ID_PATH: str = "/api/mobile/v1/urlService/url/parse"
PRODUCTS_PATH: str = "/api/mobile/v1/catalogService/catalog/search"
AUTH: dict = {
                "locationId":"50",
                "appPlatform":"WEB",
                "appVersion":0,
                "experiments":{
                    "8":"2","55":"2","58":"2","68":"1","69":"2","79":"3","99":"1","107":"2","109":"2","119":"2","120":"2","121":"2","122":"2","132":"2","144":"3","154":"1","173":"1","182":"1","184":"3","186":"2","190":"2","192":"2","194":"3","200":"2","205":"2","209":"1","218":"1","243":"1","249":"3","645":"3","646":"2","775":"2","777":"2","778":"2","790":"3","792":"3","793":"3","805":"2","808":"3","818":"2","826":"2","828":"2","837":"2","842":"2","844":"1","845":"2","852":"1","889":"1","893":"1","897":"1","899":"1","903":"1","945":"1","958":"2","962":"2","1054":"2","5779":"2","20121":"1","43568":"2","67319":"2","70070":"2","80283":"1","85160":"2","91562":"3"
                },
                "os":"UNKNOWN_OS"
}

# TODO: finish filters
class MMProductsGetter:

    def __init__(self, url):
        self.__query_url: str = url


    def __get_collection_id(self) -> str:
        json_data = {
            "url": self.__query_url[len(MM_ORIGIN):],
            "auth": AUTH
            }
        resp = requests.post(url=MM_ORIGIN+COLLECTION_ID_PATH, json=json_data, verify=True, impersonate="chrome")

        body = json.loads(resp.text)

        return body["params"]["collection"]["collectionId"]


    def GetProductsJSON(self) -> str:
        collection_id = self.__get_collection_id()

        json_data = {
            "requestVersion":12,
            "merchant":{},
            "limit":44,
            "offset":0,
            "isMultiCategorySearch":False,
            "searchByOriginalQuery":False,
            "selectedSuggestParams":[],
            "expandedFiltersIds":[],
            "sorting":0,
            "ageMore18":None,
            "showNotAvailable":True,
            "selectedAssumedCollectionId":str(collection_id),
            "selectedFilters":[],
            "collectionId":str(collection_id),
            "auth": AUTH
        }

        resp = requests.post(url=MM_ORIGIN+PRODUCTS_PATH, json=json_data, verify=True, impersonate="chrome")

        body = json.loads(resp.text)

        try:
            if not body["success"]:
                raise ValueError("error of the service response", 500)

        except Exception:
            if body["code"] == 7:
                raise Exception("the limit of the service is finished: try again later", 502)

            raise Exception("error of getting and parsed the response", 500)

        return resp.text
