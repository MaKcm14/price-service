from flask import Flask
from mm_getter import MMProductsGetter

app = Flask(__name__)

# POST-handle: client sends the {"url":"url_value"} which defines the products' resource.
# Response is the JSON-object of the products' response from mmarket service.
@app.route('/mmarket', methods=['POST']) 
def get_mmarket_products() -> str:




if __name__ == '__main__':
    app.run(host='localhost', port=8081)  # запуск работы фреймворка
