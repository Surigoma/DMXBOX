import { expect, describe, it } from "vitest";
import { render } from "vitest-browser-react";
import ErrorComponent from "./Error";

describe("Error", async () => {
    function CreateTestComponent(message: string) {
        return render(<ErrorComponent>{message}</ErrorComponent>);
    }
    it("Shown", async () => {
        const { getByTestId } = await CreateTestComponent("");
        const message = getByTestId("ErrorComponent");
        await expect.element(message).toBeVisible();
    });
    it("Shown Message", async () => {
        const { getByText } = await CreateTestComponent("message");
        const message = getByText("message");
        await expect.element(message).toBeVisible();
    });
});
