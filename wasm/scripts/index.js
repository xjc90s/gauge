define(['require', 'assert', 'wasm_exec'], (require, exports, module) => {
    require('wasm_exec');
    window.console = {
        log: (s) => document.getElementById('out').innerHTML += `> ${s}\n`,
        warn: (s) => document.getElementById('out').innerHTML += `> [WARN] ${s}\n`,
        error: (s) => document.getElementById('out').innerHTML += `> [ERROR] ${s}\n`,
    };
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("gauge.wasm"), go.importObject).then((result) => {
        go.run(result.instance);
    });

    var stepRegistry = {};
    var generalize = (s) => s.replace(/(<.*?>)|(".*?")/g, "{}");
    const caller_file = "step_implementation.js";
    window.step = (text, f) => { 
        var s = stepRegistry[generalize(text)];
        var caller_line_no = (new Error).stack.split('\n')[1].split('> Function').pop();
        var caller_line = `${caller_file}${caller_line_no}`;
        if(s){
            s.push({func: f, location: caller_line})
        } else{
            stepRegistry[generalize(text)] = [{func: f, location: caller_line}]; 
        }
    };
    window.beforeScenario = (f, opts) => {};
    window.stepImplemented = (s) => generalize(s) in stepRegistry;
    window.stepImplementationLocations = (s) => stepRegistry[s].map(i => i.location).join('|');
    window.parse = () => {
        try {
            stepRegistry = {};
            new Function(document.getElementById('implementation').innerText)();	
        } catch (error) {
            console.error(error);
        }
    };
    window.execute = (s, params) => {
        try {
            stepRegistry[s][0].func.apply({}, JSON.parse(params).map(x => x.kind === 'table' ? x.table : x.value));					
        } catch (error) {
            return error
        }
    }
});