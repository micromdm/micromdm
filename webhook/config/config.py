import os

#mysql_user = os.environ["MYSQL_USER"]
#mysql_passwd = os.environ["MYSQL_PASSWORD"]
#mysql_dbhost = os.environ["MYSQL_SERVICE_HOST"]
#mysql_dbname = os.environ["MYSQL_DATABASE"]

def const(cls):
    # Replace a class's attributes with properties,
    # and itself with an instance of its doppelganger
    is_special = lambda name: (name.startswith("__") and name.endswith("__"))
    class_contents = {n: getattr(cls, n) for n in vars(cls) if not is_special(n)}
    def unbind(value): # Get the value out of the lexical closure
        return lambda self: value
    propertified_contents = {name: property(unbind(value))
                             for (name, value) in class_contents.items()}
    receptor = type(cls.__name__, (object,), propertified_contents)
    return receptor() # Replace with an instance, so properties work

@const
class Mysql(object):
    if "MYSQL_USER" in os.environ:
        user = os.environ["MYSQL_USER"]
    else:
        user = "micromdm"

    if "MYSQL_PASSWORD" in os.environ:
        password = os.environ["MYSQL_PASSWORD"]
    else:
        password = "micromdm"

    if "MYSQL_SERVICE_HOST" in os.environ:
        service_host = os.environ["MYSQL_SERVICE_HOST"]
    else:
        service_host = "127.0.0.1"

    if "MYSQL_DATABASE" in os.environ:
        database = os.environ["MYSQL_DATABASE"]
    else:
        database = "micromdm_test"

    config = {
        'user': user,
        'password': password,
        'host': service_host,
        'database': database,
        'raise_on_warnings': True
    }

@const
class MDM(object):
    # without trailing slash!
    server_url = "https://tobis.jumpingcrab.com"

    # VPP Token downloaded from Apple. Provide path to it
    s_token = "eyJleHBEYXRlIjoiMjAxOS0xMS0wN1QwNjoyMTowMC0wODAwIiwidG9rZW4iOiJRd05hZnVnSk9vZ2tsbS9zZU52Y1pFdXhGYStLUThSNDNsQXJ4ajNQcVdUMmNpWG1iY0I5aW9LUmhOV3NSSFNEVGpIQ3NJcmxWYTAxcm5XYy9jMFNnZ09WakFQbzV2NnNSVWhVWmVjNHBMbCs3Q1pUcHZGZ0JvQm9mVTFiWTdJc3FWZnZta2UxWkZ4dVZzbUlRNEJtTnc9PSIsIm9yZ05hbWUiOiJBYmFjdXMgUmVzZWFyY2ggQUcifQ=="

    # App Id, which is used in AppStore Links
    # https://itunes.apple.com/de/app/apple-store/id1301278477?mt=8
    adam_id_str = "1301278477"
