import { useId, type ChangeEventHandler, type ReactNode } from "react";
import { NumberField as BaseNumberField } from "@base-ui/react/number-field";
import IconButton from "@mui/material/IconButton";
import FormControl from "@mui/material/FormControl";
import FormHelperText from "@mui/material/FormHelperText";
import OutlinedInput from "@mui/material/OutlinedInput";
import InputAdornment from "@mui/material/InputAdornment";
import InputLabel from "@mui/material/InputLabel";
import { MdKeyboardArrowUp } from "react-icons/md";
import { MdKeyboardArrowDown } from "react-icons/md";

/**
 * This component is a placeholder for FormControl to correctly set the shrink label state on SSR.
 */
function SSRInitialFilled(_: BaseNumberField.Root.Props) {
    return null;
}
SSRInitialFilled.muiName = "Input";

export default function NumberField({
    id: idProp,
    label,
    error,
    size = "medium",
    help,
    ...other
}: BaseNumberField.Root.Props & {
    label?: ReactNode;
    size?: "small" | "medium";
    error?: boolean;
    help?: string;
    onChange?: ChangeEventHandler;
}) {
    let id = useId();
    if (idProp) {
        id = idProp;
    }
    return (
        <BaseNumberField.Root
            {...other}
            render={(props, state) => (
                <FormControl
                    size={size}
                    ref={props.ref}
                    disabled={state.disabled}
                    required={state.required}
                    error={error}
                    variant="outlined"
                >
                    {props.children}
                </FormControl>
            )}
        >
            <SSRInitialFilled {...other} />
            <InputLabel htmlFor={id}>{label}</InputLabel>
            <BaseNumberField.Input
                id={id}
                render={(props, state) => (
                    <OutlinedInput
                        label={label}
                        inputRef={props.ref}
                        value={state.inputValue}
                        onBlur={props.onBlur}
                        onChange={props.onChange}
                        onKeyUp={props.onKeyUp}
                        onKeyDown={props.onKeyDown}
                        onFocus={props.onFocus}
                        slotProps={{
                            input: props,
                        }}
                        endAdornment={
                            <InputAdornment
                                position="end"
                                sx={{
                                    flexDirection: "column",
                                    maxHeight: "unset",
                                    alignSelf: "stretch",
                                    borderLeft: "1px solid",
                                    borderColor: "divider",
                                    ml: 0,
                                    "& button": {
                                        py: 0,
                                        flex: 1,
                                        borderRadius: 0.5,
                                    },
                                }}
                            >
                                <BaseNumberField.Increment
                                    render={
                                        <IconButton
                                            size={size}
                                            aria-label="Increase"
                                        />
                                    }
                                >
                                    <MdKeyboardArrowUp
                                        fontSize={size}
                                        style={{ transform: "translateY(2px)" }}
                                    />
                                </BaseNumberField.Increment>

                                <BaseNumberField.Decrement
                                    render={
                                        <IconButton
                                            size={size}
                                            aria-label="Decrease"
                                        />
                                    }
                                >
                                    <MdKeyboardArrowDown
                                        fontSize={size}
                                        style={{
                                            transform: "translateY(-2px)",
                                        }}
                                    />
                                </BaseNumberField.Decrement>
                            </InputAdornment>
                        }
                        sx={{ pr: 0 }}
                    />
                )}
            />
            {help !== undefined ? (
                <FormHelperText sx={{ ml: 0, "&:empty": { mt: 0 } }}>
                    {help}
                </FormHelperText>
            ) : other.min !== undefined && other.max !== undefined ? (
                <FormHelperText sx={{ ml: 0, "&:empty": { mt: 0 } }}>
                    Enter value between {other.min} and {other.max}
                </FormHelperText>
            ) : (
                <></>
            )}
        </BaseNumberField.Root>
    );
}
