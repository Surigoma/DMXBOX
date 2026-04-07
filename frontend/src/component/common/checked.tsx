import { Checkbox, FormControlLabel } from "@mui/material";
import { Controller, useFormContext } from "react-hook-form";

interface CheckedParam {
    title: string;
    target: string;
    value?: string;
}

function Checked(param: CheckedParam) {
    const { control } = useFormContext();

    function isChecked(value: any) {
        if (typeof value == "boolean") {
            return value;
        } else if (value instanceof Array) {
            return value.includes(param.value);
        }
        return false;
    }
    function onChange(value: any, checked: boolean) {
        if (typeof value == "boolean") {
            return checked;
        } else if (value instanceof Array) {
            return checked
                ? [...value, param.value]
                : value.filter((v) => v !== param.value);
        }
    }
    return (
        <Controller
            name={param.target}
            control={control}
            render={({ field }) => (
                <FormControlLabel
                    key={param.value}
                    control={
                        <Checkbox
                            checked={isChecked(field.value)}
                            onChange={(_, c) => {
                                const NEW = onChange(field.value, c);
                                field.onChange(NEW);
                            }}
                        />
                    }
                    label={param.title}
                    style={{ userSelect: "none" }}
                />
            )}
        />
    );
}

export default Checked;
