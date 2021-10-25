# Copyright 2021 magnifier Author. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
import ctypes
import io
import time

IMPL_LIB = 'libstorage.so'
_impl = ctypes.CDLL(IMPL_LIB)
_libc = ctypes.CDLL('libc.so.6')


def _str2p(s):
    return ctypes.c_char_p(bytes(s, encoding='utf-8'))


class Reader(io.IOBase):
    def __init__(self, impl) -> None:
        self.impl = impl

    def seekable(self) -> bool:
        return False

    def readable(self) -> bool:
        return True

    def writable(self) -> bool:
        return False

    def close(self) -> None:
        _impl.CloseReader(self.impl)

    def read(self, size=-1):
        if size is None:
            size = -1
        buf = ctypes.create_string_buffer(size)
        len = _impl.Read(self.impl, buf, ctypes.c_int(size))
        return buf[0:len]


class Writer(io.IOBase):
    def __init__(self, impl) -> None:
        self.impl = impl

    def seekable(self) -> bool:
        return False

    def readable(self) -> bool:
        return False

    def writable(self) -> bool:
        return True

    def close(self) -> None:
        _impl.CloseWriter(self.impl)

    def write(self, buffer):
        length = _impl.Write(
            self.impl,
            ctypes.c_char_p(buffer),
            ctypes.c_int(len(buffer)),
        )
        return length


class Stat(object):
    name: str
    mod_time: time.struct_time
    size: int
    is_dir: bool

    def __init__(self, name, mod_time, size, is_dir) -> None:
        self.name = name
        self.mod_time = mod_time
        self.size = size
        self.is_dir = is_dir

    def __str__(self):
        return f'name:{self.name}\nmod_time:{self.mod_time}\nsize:{self.size}\nis_dir:{self.is_dir}'


class Storage(object):
    def __init__(self, configs):
        self.driver = _impl.CreateStorageDriver(
            _str2p('filesystem'),
            _str2p(configs),
        )

    def get_content(self, path):
        p = _impl.GetContent(self.driver, _str2p(path))
        result = ctypes.string_at(p)
        _libc.free(p)
        return result

    def put_content(self, path, data):
        _impl.PutContent(
            self.driver,
            _str2p(path),
            ctypes.c_char_p(data),
            ctypes.c_int(len(data)),
        )

    def reader(self, path, offset):
        impl = _impl.Reader(
            self.driver,
            _str2p(path),
            ctypes.c_int(offset),
        )
        return Reader(impl)

    def writer(self, path, append):
        impl = _impl.Writer(
            self.driver,
            _str2p(path),
            ctypes.c_bool(append)
        )
        return Writer(impl)

    def stat(self, path):
        p = _impl.Stat(self.driver, _str2p(path))
        s = ctypes.string_at(p)
        _libc.free(p)
        name, mod, size, is_dir = str(s, encoding='utf-8').split(':')
        return Stat(
            name,
            time.localtime(int(mod, 10)),
            int(size, 10),
            bool(is_dir),
        )

    def list(self, path):
        p = _impl.List(self.driver, _str2p(path))
        result = ctypes.string_at(p)
        _libc.free(p)
        return str(result, encoding='utf-8').split('\n')

    def move(self, src, dst):
        _impl.Move(self.driver, _str2p(src), _str2p(dst))

    def delete(self, path):
        _impl.Delete(self.driver, _str2p(path))

    def close(self):
        _impl.ReleaseDriver(self.driver)
