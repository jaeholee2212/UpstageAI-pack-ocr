from enum import Enum
from typing import Dict, Optional
import json


class FieldType(Enum):
    INT = 1
    STR = 2
    SET = 3


class Snorkel:

    def __init__(self, table: str, token: str) -> None:
        self.table = table
        self.token = token
        self.spec = {}

    def add_int_field(self, name: str):
        self.__add_field(name, FieldType.INT)

    def add_str_field(self, name: str):
        self.__add_field(name, FieldType.STR)

    def add_set_field(self, name: str):
        self.__add_field(name, FieldType.SET)

    def __add_field(self, name: str, ftype: FieldType):
        self.spec[name] = ftype

    def write(self, data: dict[str, int | str | list[str]]):
        out = {}
        for k, v in data.items():
            if not k in self.spec:
                continue

            ftype = self.spec[k]
            if ftype == FieldType.INT:
                intv = Snorkel.to_int(v)
                if intv is None:
                    raise ValueError(
                        f'a value is not compatible with int, k={k}')
                out[k] = intv
            elif ftype == FieldType.STR:
                strv = f'{v}'
                out[k] = strv
            elif ftype == FieldType.SET:
                if not isinstance(v, list):
                    raise ValueError(f'expected a list of string, k={k}')
                out[k] = [f'{x}' for x in list]
            else:
                raise ValueError(f'an unknown type is given, k={k}')

        evt = {'table': self.table, 'token': self.token, 'data': out}
        evt_json = json.dumps(evt)
        print(f'{evt_json}\n')

    @staticmethod
    def to_int(v: int | str | list[str]) -> int | None:
        t = type(v)
        if t == int or t == float:
            return int(v)
        if t == str:
            try:
                return int(v)
            except Exception:
                return None

        return None
