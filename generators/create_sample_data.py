# -*- coding: utf-8 -*-
"""
Created on Sat Oct  2 10:59:36 2021
@author: Markus Horst
"""
#%% imports and initialization
from faker import Faker


import json
import pandas as pd
import numpy as np
import random
import datetime
from scipy.stats import geom
from tqdm import tqdm

some_seed = 42
np.random.seed(some_seed)
n_adresses = 100000
n_persons  = 10000
n_invoices = 1000
min_net_sum = 20

# create hierachy with dictionary   a: (b,c)  a is level, (repeated) sampling of integer n:  b <= n <= c as number of agents in that level
agent_hierarchy = {1:(1,1), 2:(5,5), 3:(2,5), 4:(1,5), 5:(0,5)}  
#agent_hierarchy = {1:(1,1), 2:(5,5), 3:(5,5), 4:(5,5), 5:(5,5)}  

def jsonobj_hook(obj):
    value = obj.get("features")
    pbar = None
    if value:
        pbar = tqdm(value)
    if pbar is not None:
        for item in pbar:
            pass
            pbar.set_description("Loading")
    return obj


#%% read geojson-files with adress data from Germany (takes some minutes)
create_adress_store = False
if create_adress_store:
    path_to_file = r'./input_data/squares.geojson'
    with open(path_to_file, encoding='utf-8') as f:
       strasse_ort_01 = json.load(f, object_hook=jsonobj_hook)
    df_squares = pd.json_normalize(strasse_ort_01['features'])
    df_squares = df_squares[ df_squares['properties.HIGHWAY'].str.contains('residential')   |
                             df_squares['properties.HIGHWAY'].str.contains('living_street') |
                             df_squares['properties.HIGHWAY'].str.contains('pedestrian')     ]
    df_squares = df_squares[ df_squares['properties.NAME'].str.len().between(2,50)]


    path_to_file = r'./input_data/streets.geojson'
    with open(path_to_file, encoding='utf-8') as f:
        strasse_ort_02 = json.load(f, object_hook=jsonobj_hook)
    df_streets = pd.json_normalize(strasse_ort_02['features'])
    df_streets = df_streets[ df_streets['properties.HIGHWAY'].str.contains('residential')   |
                             df_streets['properties.HIGHWAY'].str.contains('living_street') |
                             df_streets['properties.HIGHWAY'].str.contains('pedestrian')     ]
    df_streets = df_streets[ df_streets['properties.NAME'].str.len().between(2,50)]

    df_streets_and_squares = df_streets.append(df_squares, ignore_index=True)
    df_streets_and_squares = df_streets_and_squares[['properties.NAME', 'properties.GEMEINDE', 'properties.PLZ']]
    df_streets_and_squares.rename(columns={'properties.NAME': 'Street', 
                                           'properties.GEMEINDE': 'Residence', 'properties.PLZ': 'ZipCode'}, inplace = True)
    df_streets_and_squares['ZipCode'] = df_streets_and_squares.apply(lambda x: x['ZipCode'][0:(x['ZipCode'].find(','))] if ',' in x['ZipCode'] else x['ZipCode'], axis = 1)
    df_streets_and_squares['ZipCode'] = df_streets_and_squares['ZipCode'].apply(lambda zc: int(zc))
    df_streets_and_squares.drop_duplicates( keep='first', inplace = True)

    df_plz_vorwahl = pd.read_csv(r"./input_data/orte_de_plz_vorwahl.csv", 
                                 sep = ';', usecols=['Plz', 'Vorwahl'], dtype={'Plz': np.int64, 'Vorwahl': str})
    df_plz_vorwahl.drop_duplicates( subset = 'Plz', keep='first', inplace = True)
    df_plz_vorwahl.dropna(inplace=True)

    df_streets_and_squares_vorwahl = pd.merge(df_streets_and_squares, df_plz_vorwahl, left_on=['ZipCode'], right_on=['Plz'], how='inner')
    df_streets_and_squares_vorwahl = df_streets_and_squares_vorwahl.drop('Plz', axis=1)
    df_sample_adresses = df_streets_and_squares_vorwahl.sample(n=n_adresses)
    
    out_json = df_sample_adresses.to_json(orient='records')
    with open(r'.\adresses.json', 'w') as f:
        f.write(out_json)
 
#%% load adresses

df_adress_store = pd.read_json(r'./adresses.json', 
                               orient='records',  dtype = {"ZipCode": object, "Vorwahl": object})

#%% load roles
roles = pd.read_json("input_data/roles.json")
print(roles)
rolecount = len(roles.transpose().index)
#%% create persons

def get_role_dict_by_id(id):
    return roles[id]

house_numbers = geom.rvs(0.03, size = n_persons) 
ascii_lowercase = [chr(i) for i in range(ord('a'), ord('z') + 1)]
postfixes =[' '] + ascii_lowercase
df_sample_adresses = df_adress_store.sample(n=n_persons, replace=True)
house_numbers_postfix_index =(geom.rvs(0.9, loc=-1, size = n_persons))

df_first_names_m = pd.read_csv(r"./input_data/vornamen_m.txt", names=['FirstName'], skiprows = 1)
df_first_names_f = pd.read_csv(r"./input_data/vornamen_w.txt", names=['FirstName'], skiprows = 1)
df_last_names = pd.read_csv(r"./input_data/nachnamen.txt", names=['Name'], skiprows = 1)
# some first names have  similar spellings of the name, separated by '/'. Use only the first one
df_first_names_m.iloc[:,0] = df_first_names_m.apply(lambda x: x[0][0:(x[0].find('/'))] if '/' in x[0] else x[0], axis = 1)
df_first_names_f.iloc[:,0] = df_first_names_f.apply(lambda x: x[0][0:(x[0].find('/'))] if '/' in x[0] else x[0], axis = 1)

n_persons_f = int(0.51*n_persons)
n_persons_m = n_persons - n_persons_f

df_persons = df_last_names.sample(n=n_persons, replace=True, weights=geom.pmf(p=0.005, k = np.arange(len(df_last_names.index))))
df_persons.reset_index(drop=True, inplace=True) 
df_first_names = df_first_names_f.sample(n=n_persons_f, replace=True, weights=geom.pmf(p=0.01, 
  k = np.arange(len(df_first_names_f.index)))).append(df_first_names_m.sample(n=n_persons_m, replace=True, weights=geom.pmf(p=0.01, 
  k = np.arange(len(df_first_names_m.index)))),ignore_index=True)

df_persons['FirstName']= df_first_names.iloc[:,0]
df_persons['Street']= df_sample_adresses.loc[:,'Street'].values
df_persons['Residence']= df_sample_adresses.loc[:,'Residence'].values
df_persons['ZipCode']= df_sample_adresses.loc[:,'ZipCode'].values
df_persons['PhoneNumber_pre']= df_sample_adresses.loc[:,'Vorwahl'].values
df_persons['HouseNumber_onlyNumber']= house_numbers
df_persons['house_numbers_postfix_index'] = house_numbers_postfix_index
df_persons['house_numbers_postfix'] = df_persons.apply(lambda x: postfixes[x['house_numbers_postfix_index']].strip() , axis = 1)
df_persons['HouseNumber']= df_persons['HouseNumber_onlyNumber'].astype(str) + df_persons['house_numbers_postfix'] 

start_date = datetime.date(1940, 1, 1)
end_date = datetime.date(1980, 1, 1)
time_between_dates = end_date - start_date
days_between_dates = time_between_dates.days
df_persons['BirthDate'] = df_persons.apply(lambda x: datetime.datetime.combine(start_date + datetime.timedelta(days=random.randrange(days_between_dates)), datetime.datetime.min.time()).isoformat() + "Z", axis = 1)

start_date = datetime.date(2000, 1, 1)
end_date = datetime.date(2009, 12, 31)
time_between_dates = end_date - start_date
days_between_dates = time_between_dates.days
df_persons['RegistrationDate'] = df_persons.apply(lambda x: datetime.datetime.combine(start_date + datetime.timedelta(days=random.randrange(days_between_dates)), datetime.datetime.min.time()).isoformat() + "Z" , axis = 1)
df_persons['PhoneNumber'] = df_persons.apply(lambda x: str(x['PhoneNumber_pre']) +' '+ str(random.randint(1000, 9999999)), axis = 1)

locale_list = ['de_DE']
fake = Faker(locale_list)
Faker.seed(some_seed)
df_persons['EmailAddress'] = df_persons.apply(lambda x: fake.email(), axis = 1)

df_persons['PersonID'] = df_persons.index + 1
df_persons['Role'] = [get_role_dict_by_id(0).to_dict() for x in range(n_persons)]
df_persons['RoleID'] = 1
df_persons.drop(['PhoneNumber_pre', 'house_numbers_postfix_index', 'HouseNumber_onlyNumber', 'house_numbers_postfix'],axis=1,inplace=True)
#after the drop of unused cols df_persons is complete


#%% create hierarchy of agents
np.random.seed(some_seed)

agent_hierarchy_n = {}
n_agents = 0
n_draws = 1
for level in agent_hierarchy:
    (min_n, max_n) = agent_hierarchy[level]
    n_same_level_agents = 0
    agent_hierarchy_n[level] = []
    for i in tqdm(range(0, n_draws)):
        if min_n < max_n:
            n_actual_agents = np.random.randint(min_n, max_n)
        else:
            n_actual_agents = min_n
        n_same_level_agents += n_actual_agents
        n_agents += n_actual_agents
        agent_hierarchy_n[level].append(n_actual_agents)
    print(level, n_same_level_agents, n_agents)
    n_draws = n_same_level_agents

df_agents = df_persons.sample(n_agents, replace = False)
#
pos = 0
lst_agents = []
lst_agent_ids = []
lst_supervisor_ids = []
for level in agent_hierarchy_n:
    l = level
    supervisor_offset = 0
    for agent_nr in agent_hierarchy_n[level]:
        if l == 1:
            lst_supervisors = [-1] * agent_nr
            lst_supervisor_ids = [-1] * agent_nr
        for j in range(agent_nr):
            if rolecount <= level:
                l = rolecount - 1
                print("NOTICE: moving lower level persons to last available level - please define other levels in input/roles.json if you wish to expand levels.")

            lst_agents.append(df_agents.iloc[pos].to_dict())
            person = df_persons.iloc[lst_agents[-1]["PersonID"]-1]
            person.Role = get_role_dict_by_id(l).to_dict()
            person.RoleID = person.Role["RoleID"]
            df_persons.at[lst_agents[-1]["PersonID"]-1] = person
            df_agents.iloc[pos] = person
            lst_agent_ids.append(df_agents.iloc[pos].PersonID)
            lst_agents[len(lst_agents) -1] = df_agents.iloc[pos].to_dict()
            pos += 1
            if level < list(agent_hierarchy_n.keys())[-1] :
                lst_supervisors.extend( [lst_agents[-1]] * agent_hierarchy_n[level+1][j + supervisor_offset]  )
        supervisor_offset = supervisor_offset + agent_nr
modification_date = datetime.datetime(2021, 1, 1).isoformat() + "Z"
df_hierarchy = pd.DataFrame({'Agent': lst_agents, 'AgentID': lst_agent_ids, 'Supervisor': lst_supervisors, 
                             'ModificationDate': modification_date, 'AgentStatus': 1})
df_hierarchy.loc[ df_hierarchy["Supervisor"]== -1 , "Supervisor"] = None

def getSupervisorID(supervisor):
    if supervisor is not None:
        return supervisor["PersonID"]
    else:
        return None
df_hierarchy["SupervisorID"] = [getSupervisorID(supervisor) for supervisor in df_hierarchy["Supervisor"]]
df_hierarchy["SupervisorID"] = df_hierarchy["SupervisorID"].astype(pd.Int64Dtype())


#%% create invoices



np.random.seed(some_seed)

df_invoices = pd.DataFrame({'NetSum': np.random.normal(150, 30, size = n_invoices)}).round(2)
df_invoices.loc[ df_invoices.NetSum < min_net_sum, 'NetSum' ] = min_net_sum
df_invoices.round(2)
df_invoices['GrossSum'] = (df_invoices['NetSum'] * 1.19).round(2)
df_invoices['VAT'] = (df_invoices['GrossSum'] - df_invoices['NetSum']).round(2)


start_date = datetime.date(2010, 1, 1)
end_date = datetime.date(2020, 12, 31)

time_between_dates = end_date - start_date
days_between_dates = time_between_dates.days
df_invoices['InvoiceDate'] = df_invoices.apply(lambda x: (start_date + datetime.timedelta(days=random.randrange(days_between_dates))) , axis = 1)
df_invoices['PayDate']  = df_invoices['InvoiceDate'].apply(lambda day: datetime.datetime.combine(day + datetime.timedelta(days=10), datetime.datetime.min.time()).isoformat() + "Z")
df_invoices['InvoiceDate'] = df_invoices['InvoiceDate'].apply(lambda day: datetime.datetime.combine(day, datetime.datetime.min.time()).isoformat() + "Z")
df_invoices['OpenSum'] = 0
df_invoices['Customer'] = list(df_persons[df_persons.RoleID==1].sample(n_invoices, replace = True).to_dict('records'))
df_invoices['Agent'] =  list(df_persons[df_persons.RoleID!=1].sample(n_invoices, replace = True).to_dict('records'))
df_invoices['AgentID'] = [person["PersonID"] for person in df_invoices.Agent]
df_invoices['InvoiceID'] = df_invoices.index + 1


def all_but(*names, df):
    names = set(names)
    res = [item for item in df if item not in names]
    return res
#%% write json-files
out_json = (df_persons
    .to_json(orient='records'))
with open(r'./output_data/persons.json', 'w') as f:
    f.write(out_json)
    
out_json = df_hierarchy.to_json(orient='records')
# out_json = df_hierarchy.apply(lambda x: [x.dropna()], axis=1).to_json()
with open(r'./output_data/hierarchy.json', 'w') as f:
    f.write(out_json)  
    
out_json = df_invoices.to_json(orient='records')
out_json = out_json.replace('{"NetSum"', '\n{"NetSum"')
with open(r'./output_data/invoices.json', 'w') as f:
    f.write(out_json)