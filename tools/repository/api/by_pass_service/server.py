from flask import Flask, request
from megamarket import MegaMarketAPI

app = Flask(__name__)

# POST-handle: client sends the next JSON-object:
# {
#   "query": "query_text",
#   "sample": "sample_num",
#   "sort": "sort_num",
#   "show_not_available":flag,
#   "is_price_filter_set": flag_filter_set
#   "price_filter": {
#       "price_down": "low_price_border",
#       "price_up": "high_price_border",
#   }
# }
# which defines the products' request query.
# Response is the JSON-object of the response from megamarket service.
@app.route('/mmarket', methods=['POST'])
def get_mmarket_products():
    body = request.get_json()
    try:
        getter = MegaMarketAPI(body)
        return getter.get_products_json(), 200

    except AttributeError | TypeError:
        return str({"error": "error of the request's structure"}), 400

    except OverflowError as excp:
        return str({"error": excp.args[0]}), 502

    except Exception as excp:
        return str({"error": excp.args[0]}), 500


if __name__ == '__main__':
    app.run(host='localhost', port=8081)
