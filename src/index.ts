import {info, getInput, setFailed, setOutput, setSecret} from "@actions/core";
import {AppTokenService} from "./services/app-token-service";

(async () => {
    const appId = getInput("app_id", {required: true});
    console.log(`appId: ${appId}`);
    const privateKey = getInput("private_key", {required: true});
    console.log(`privateKey: ${privateKey}`);
    const installationId = getInput("installation_id");
    console.log(`installationId: ${installationId}`);
    const repository = getInput("repository");
    console.log(`repository: ${repository}`);
    const appTokenSvc = new AppTokenService({
        appId,
        privateKey,
        installationId,
        repository
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
