import pandas as pd
import numpy as np
import streamlit as st
from streamlit.runtime.caching.cache_data_api import CacheDataAPI
import sqlalchemy as sa
from common import query, get_conn, has_scopes

st.set_page_config('Admin')

if not has_scopes(['admin']):
    st.error('Unauthorized')
    st.stop()

@st.cache_data
def get_users():
    q = sa.text('''
    select u.email from auth.users u order by u.email
    ''')
    return query(q)

@st.cache_data
def get_all_scopes():
    q = sa.text('select s.key from auth.scopes s')
    return query(q)['key']

def get_user_scopes(user: str):
    q = sa.text('''
    SELECT s.key
    FROM auth.users_to_scopes us
    JOIN (SELECT * FROM auth.users WHERE email = :email) u ON u.id = us.user_id
    JOIN auth.scopes s ON s.id = us.scope_id;
    ''').bindparams(email=user)
    return query(q)['key']

def set_user_scopes(user: str, scopes: list[str]):
    q = sa.text('''SELECT auth.user_set_scopes(:user, :scopes)
    ''').bindparams(user=user, scopes=scopes)
    with get_conn() as c:
        c.execute(q)
        c.commit()

def refresh_scopes():
    '''
    Refresh cached data and reset list of scopes
    '''
    st.session_state.admin_scope_edit = list(get_user_scopes(sel_user))

def delete_user(user):
    q = sa.text('''
    DELETE FROM auth.users WHERE email = :email
    ''').bindparams(email=user)
    with get_conn() as c:
        c.execute(q)
        c.commit()
    st.cache_data.clear()
    st.session_state.pop('user_email')
    st.session_state['delete_user_cb'] = False

def add_user(user):
    q = sa.text('''
    INSERT INTO auth.users (email) VALUES (:email)
    RETURNING users.id;
    ''').bindparams(email=user)
    with get_conn() as c:
        c.execute(q)
        c.commit()
    st.cache_data.clear()
    st.session_state['user_email'] = user

sel_user = st.selectbox('User', get_users().email, key='user_email')
user_scopes = get_user_scopes(sel_user)
scope_edit = st.multiselect(
    'Scopes',
    get_all_scopes(),
    default=user_scopes,
    key='admin_scope_edit')

with st.empty():
    c1, c2, _ = st.columns(3)
    c1.button('Commit', on_click=lambda: set_user_scopes(sel_user, list(scope_edit)))
    c2.button('Refresh', on_click=refresh_scopes)

with st.empty():
    cBtn, cConfirm, _ = st.columns(3)
    if cBtn.checkbox('Delete User', key='delete_user_cb'):
        cConfirm.button('Are you sure?', on_click=lambda: delete_user(sel_user))
        
with st.expander('Add User'):
    new_name = st.text_input('Email')
    st.button(
        'Add User',
        disabled=(new_name == ''),
        key='add_user_btn',
        on_click=lambda: add_user(new_name),
    )
