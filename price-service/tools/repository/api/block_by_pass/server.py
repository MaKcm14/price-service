from flask import Flask, request
from mm_getter import MMProductsGetter

app = Flask(__name__)

# POST-handle: client sends the {"query":"query_text"} 
# which defines the products' request query.
# Response is the JSON-object of the products response from megamarket service.
@app.route('/mmarket', methods=['POST']) 
def get_mmarket_products():
    body = request.get_json()
    try:
        getter = MMProductsGetter(body["query"])
        return getter.get_products_json(), 200

    except AttributeError | TypeError:
        return str({"error": "error of the request's structure"}), 400

    except OverflowError as excp:
        return str({"error": excp.args[0]}), 502

    except Exception as excp:
        return str({"error": excp.args[0]}), 500


if __name__ == '__main__':
    app.run(host='localhost', port=8081)
