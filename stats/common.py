import pandas as pd
import streamlit as st
import streamlit.web.server.websocket_headers as wsheaders
from dotenv import load_dotenv
import sqlalchemy as sa
import os

# For development
load_dotenv('../.env')

@st.cache_resource
def get_db():
    driver = 'postgresql+psycopg2'
    user = os.getenv('POSTGRES_USER')
    passw = os.getenv('POSTGRES_PASSWORD')
    host = os.getenv('POSTGRES_HOST')
    engine = sa.create_engine(f'{driver}://{user}:{passw}@{host}/postgres')
    return engine

def query(q: sa.TextClause) -> pd.DataFrame:
    with get_db().connect() as c:
        return pd.read_sql(q, c)
    
def get_conn() -> sa.Connection:
    return get_db().connect()

def has_scopes(scopes: list[str]) -> bool:
    if len(scopes) == 0:
        return True
    headers = wsheaders._get_websocket_headers() or {}
    email = headers.get('X-Forwarded-Email')
    if email is None:
        return False
    q = sa.text('SELECT auth.user_has_scopes(:email, :scopes)').bindparams(email=email, scopes=scopes)
    with get_db().connect() as c:
        return c.scalar(q)
