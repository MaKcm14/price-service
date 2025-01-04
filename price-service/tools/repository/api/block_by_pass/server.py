from flask import Flask, request, jsonify
from mm_getter import MMProductsGetter

app = Flask(__name__)

# POST-handle: client sends the {"url":"url_value"} which defines the products' resource.
# Response is the JSON-object of the products' response from mmarket service.
@app.route('/mmarket', methods=['POST']) 
def get_mmarket_products():
    body = request.get_json()
    try:
        getter = MMProductsGetter(body["url"])
        return getter.GetProductsJSON(), 200

    except OverflowError as excp:
        return str({"error": excp[0]}), 502

    except Exception as excp:
        return str({"error": excp[0]}), 500




if __name__ == '__main__':
    app.run(host='localhost', port=8081)
