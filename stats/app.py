import altair as alt
import math
import pandas as pd
import numpy as np
import streamlit as st
from dotenv import load_dotenv
import sqlalchemy as sa
import os


st.set_page_config(layout='wide')

# Rows per page for paged queries
PAGE_SIZE = 40

# For development
load_dotenv('../.env')

"""
# sts-metrics Statistics Views
Made with Streamlit!

Links: [documentation](https://docs.streamlit.io) and [community forums](https://discuss.streamlit.io).
"""

@st.cache_resource
def db_conn():
    driver = 'postgresql+psycopg2'
    user = os.getenv('POSTGRES_USER')
    passw = os.getenv('POSTGRES_PASSWORD')
    host = os.getenv('POSTGRES_HOST')
    engine = sa.create_engine(f'{driver}://{user}:{passw}@{host}/postgres')
    return engine.connect()

def query(sql, params=[]):
    return pd.read_sql(sa.text(sql), db_conn(), params=params)


@st.cache_data
def char_map() -> dict[str,int]:
    '''List of character names'''
    df = query('select c.* from character_list c')
    return dict(zip(df.name, df.id))

@st.cache_data
def extra_data_row_count(character: str) -> int:
    '''Number of runs_extra rows the character has'''
    char_id = char_map()[character]
    q = sa.text('''
    with run_ids as (select r.id from runsdata r where r.character_id = :char_id)
    select count(e.run_id) from runs_extra e
    inner join run_ids r on r.id = e.run_id
    ''').bindparams(char_id=char_id)
    v = db_conn().scalar(q)
    print(v)
    return v

def pg_array(ar: list[any]) -> str:
    return '{'+','.join(ar)+'}'


def quartiles_to_pd(quarts: pd.DataFrame, prefix: str) -> pd.DataFrame:
    '''
    Helper function that takes a pandas column of 3-tuples, and converts it into
    a DataFrame with 3 columns, prefixed with 'prefix'.
    '''
    cols = [prefix + s + ' ' for s in ['Q25', 'Q50', 'Q75']]
    data = [np.array(x) for x in quarts.to_numpy()]
    return pd.DataFrame(data=data, columns=cols)


tabOverview, tabCharacter = st.tabs(['Overview', 'Character'])

with tabOverview:
    st.text('''
    Reset the database connection and some other cached information.
    ''')
    if st.button('Clear Cache'):
        st.cache_data.clear()
        st.cache_resource.clear()

    def view_overview():
        st.header('Overview')
        char_list = char_map().keys()
        chars = st.multiselect('Character(s)', char_list, default=char_list)
        q = sa.text('''
        select s.name, s.runs, s.wins, s.avg_win_rate, s.p_deck_size, s.p_floor_reached
        from stats_overview s
        where s.name = any(:p1 :\:text[])
        ''').bindparams(p1=chars)
        st.table(pd.read_sql(q, db_conn()))
    view_overview()

    def view_build_versions():
        '''
        ## List of known build versions
        '''
        st.markdown(view_build_versions.__doc__)
        q = sa.text('''
        select distinct s.str as version from strcache s
        join runsdata r on s.id = r.build_version
        ''')
        st.table(pd.read_sql(q, db_conn()))
    view_build_versions()


with tabCharacter:
    # For all per-character 
    st.markdown('## Character Choice')
    character = st.selectbox('Character', char_map().keys())

    def view_card_counts():
        '''
        ## Card Counts
        '''
        st.markdown(view_card_counts.__doc__)
        q = sa.text('''
        with char_id as (select S.id from strcache S where S.str = :p1),
        q1 as (
            select S.str as name, ST.total, ST.upgrades from stats_card_counts ST
            join strcache S on ST.card_id = S.id
            join char_id on ST.char_id = char_id.id
        )
        select * from q1 where total > 1 order by total desc
        ''').bindparams(p1=character)
        st.dataframe(pd.read_sql(q, db_conn()))
    view_card_counts()

    @st.cache_data
    def card_stats(character):
        char_id = char_map()[character]
        q = sa.text('SELECT * FROM per_character_card_stats(:char_id)')\
            .bindparams(char_id=char_id)
        return pd.read_sql(q, db_conn())

    def view_card_quartiles():
        '''
        ## Card Statistics
        Q25, Q50, Q75 are the quartiles amongst all the data. Showing the quartiles
        can tell you a lot more about what's happening than just averages.
        '''
        KEYP = view_card_quartiles.__name__
        st.markdown(view_card_quartiles.__doc__)
        df = card_stats(character)
        df_deck = quartiles_to_pd(df['deck'], 'Deck ')
        df_floors = quartiles_to_pd(df['floor'], 'Floor ')
        df2 = df[['card', 'runs', 'wins']].copy().join([df_deck, df_floors])
        st.dataframe(df2)
    view_card_quartiles()

    def view_extra_data():
        '''
        ## Additional Un-Parsed JSON Fields
        '''
        KEYP = view_extra_data.__name__
        st.markdown(view_extra_data.__doc__)
        page_max = math.ceil(extra_data_row_count(character) / PAGE_SIZE)
        page_num = st.number_input('Page', 1, page_max, step=1, key=KEYP+'page')
        char_id = char_map()[character]
        offset = PAGE_SIZE * (page_num - 1)
        q = sa.text('''
        with run_ids as (select r.id from runsdata r where r.character_id = :char_id)
        select e.run_id, e.extra from runs_extra e
        inner join run_ids r on r.id = e.run_id
        order by e.run_id offset :offset limit :page_size
        ''').bindparams(char_id=char_id, offset=offset, page_size=PAGE_SIZE)
        st.dataframe(pd.read_sql(q, db_conn()))
    view_extra_data()