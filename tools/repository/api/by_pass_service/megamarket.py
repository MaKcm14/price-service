import json
from curl_cffi import requests

MM_ORIGIN: str = "https://megamarket.ru"
PRODUCTS_PATH: str = "/api/mobile/v1/catalogService/catalog/search"
FILTER_ID = "88C83F68482F447C9F4E401955196697"
AUTH: dict = {
    "locationId":"50",
    "appPlatform":"WEB",
    "appVersion":0,
    "experiments":{
        "8":"2",
        "55":"2",
        "58":"2",
        "68":"1",
        "69":"2",
        "79":"3",
        "99":"1",
        "107":"2",
        "109":"2",
        "119":"2",
        "120":"2",
        "121":"2",
        "122":"2",
        "132":"2",
        "144":"3",
        "154":"1",
        "173":"1",
        "182":"1",
        "184":"3",
        "186":"2",
        "190":"2",
        "192":"2",
        "194":"3",
        "200":"2",
        "205":"2",
        "209":"1",
        "218":"1",
        "243":"1",
        "249":"3",
        "645":"3",
        "646":"2",
        "775":"2",
        "777":"2",
        "778":"2",
        "790":"3",
        "792":"3",
        "793":"3",
        "805":"2",
        "808":"3",
        "818":"2",
        "826":"2",
        "828":"2",
        "837":"2",
        "842":"2",
        "844":"1",
        "845":"2",
        "852":"1",
        "889":"1",
        "893":"1",
        "897":"1",
        "899":"1",
        "903":"1",
        "945":"1",
        "958":"2",
        "962":"2",
        "1054":"2",
        "5779":"2",
        "20121":"1",
        "43568":"2",
        "67319":"2",
        "70070":"2",
        "80283":"1",
        "85160":"2",
        "91562":"3"
    },
    "os":"UNKNOWN_OS"
}

class MegaMarketAPI:
    __err_response_struct: str = "mmarket: error of the response's structure"
    __err_service_limit: str = "mmarket: the limit of the service is over: try again later"
    __err_service_interaction: str = "mmarket: error of the service interaction"


    def __init__(self, body):
        self.__query: str = body["query"]
        self.__page: int = int(body["sample"])
        self.__sort: int = int(body["sort"])
        self.__show_not_available: bool = body["show_not_available"]
        self.__flag_price_filter: bool = body["is_price_filter_set"]
        
        if self.__flag_price_filter:
            self.__price_range: tuple = (
                body["price_filter"]["price_down"], body["price_filter"]["price_up"])


    def __send_request(self, json_data):
        resp = requests.post(url=MM_ORIGIN+PRODUCTS_PATH, json=json_data, verify=True, impersonate="chrome")
        body = json.loads(resp.text)

        try:
            if not body["success"]:
                raise BaseException(self.__err_service_interaction)

        except AttributeError | TypeError:
            raise BaseException(self.__err_response_struct)

        except BaseException:
            if body["code"] == 7:
                raise OverflowError(self.__err_service_limit)
            raise

        return resp


    def get_products_json(self) -> str:
        json_data = {
            "requestVersion" : 12,
            "merchant" : {},
            "limit" : 44,
            "offset" : (self.__page - 1) * 44,
            "isMultiCategorySearch" : False,
            "searchByOriginalQuery" : False,
            "selectedSuggestParams" : [],
            "expandedFiltersIds" : [],
            "sorting": self.__sort,
            "ageMore18" : None,
            "showNotAvailable" : self.__show_not_available,
            "selectedFilters" : [],
            "searchText" : self.__query,
            "auth" : AUTH
        }

        if self.__flag_price_filter:
            json_data["selectedFilters"] = [{
                "filterId" : FILTER_ID,
                "type" : 1,
                "value" : self.__price_range[0]
            },
            {
                "filterId" : FILTER_ID,
                "type" : 2,
                "value" : self.__price_range[1]
            },
            ]

        resp = self.__send_request(json_data)
        body = json.loads(resp.text)

        try:
            if len(body["items"]) == 0:
                if not body["success"]:
                    raise BaseException(self.__err_service_interaction)

                collection_id = body["processor"]["collectionId"]

                json_data["collectionId"] = collection_id
                json_data["selectedAssumedCollectionId"] = collection_id
                json_data["requestVersion"] = 10
                json_data.pop("searchText")

                resp = self.__send_request(json_data)
    
        except AttributeError | TypeError:
            raise BaseException(self.__err_response_struct)

        return resp.text

