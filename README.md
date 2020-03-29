# Luna - Wrapping the Moon

More README coming soon, this project is being extracted from `dragon-mud` with
some syntax improvements/cleanup and comment improvements.

Becuase of that, still very much a WIP, the API will change between versions. And
this README will grow to accomodate.

## Special Thanks

Luna is built on top of the very excellent [gopher-lua](https://github.com/yuin/gopher-lua)
to provide a Lua state for executing source code inside of your Go library. If you're looking
for something more lightweight than Luna then you should look no further than here.

Luna is also powered by a lightweight wrapper for gopher-lua, [gopher-luar](https://github.com/layeh/gopher-luar)
which uses reflection to enable more flexibility in passing data to and from a Lua state.
