import {info, setFailed, setOutput, setSecret} from "@actions/core";
import {AppTokenService} from "./services/app-token-service";

(async () => {
    const appTokenSvc = new AppTokenService();

    try {
        const appToken = await appTokenSvc.getToken();

        setSecret(appToken);
        setOutput("app_token", appToken);
        info("Token generated successfully: 🔑");
    } catch (e) {
        setFailed(e as unknown as string | Error)
    }
})();
