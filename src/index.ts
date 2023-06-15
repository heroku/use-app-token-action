import {info, getInput, setFailed, setOutput, setSecret} from "@actions/core";
import {AppTokenService} from "./services/app-token-service";

(async () => {
    const appTokenSvc = new AppTokenService({
        appId: getInput("app_id", {required: true}),
        privateKey: getInput("private_key", {required: true}),
        installationId: getInput("installation_id"),
        repository: getInput("repository")
    });

    try {
        console.info("Starting execution: Use GitHub App Token Action");

        const appToken = await appTokenSvc.getToken();

        setSecret(appToken);
        setOutput("app_token", appToken);
        info("Token generated successfully: 🔑");
    } catch (e) {
        setFailed(e as unknown as string | Error)
    }
})();
