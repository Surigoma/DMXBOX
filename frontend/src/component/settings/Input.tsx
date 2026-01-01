import type { Config } from "../../types";
import Modules from "./Modules";

interface InputsParams {
    config: Config;
}
function Inputs(param: InputsParams) {
    return <Modules config={param.config.modules}></Modules>;
}

export default Inputs;
