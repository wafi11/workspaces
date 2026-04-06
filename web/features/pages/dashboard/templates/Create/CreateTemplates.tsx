"use client"

import { STEPS, useCreateTemplate } from "@/features/hooks/templates/useCreateTemplate";
import { StepBar } from "./components/StepBar";
import { StepBasicInfo } from "./components/StepBasicInfo";
import { StepVariables } from "./components/StepVariables";
import { StepAddons } from "./components/StepAddons";
import { StepReview } from "./components/StepReview";
import { Button } from "@/components/ui/button";
import { TopBarAdmin } from "@/components/layouts/TopBarAdmin";
import { StepFiles } from "./components/StepFiles";


export default function CreateTemplatePage(){
    const hook = useCreateTemplate();
    
    
    
    return (
        <>
        <TopBarAdmin title="Form Templates" description="Form New Templates"/>
                <StepBar steps={STEPS} current={hook.step} onGoTo={hook.goTo} />
                {StepForm({
                    hook : hook
                })}
                <div className="flex justify-end mt-4 gap-3 ">
                {
                    hook.step !== 0 && (
                        <Button onClick={() =>hook.prev()}>
                            Previous
                        </Button>
                    )
                }
                {
                    hook.step !== 4 ? (
                        <Button onClick={() =>hook.next()}>
                            Next
                        </Button>
                    ) : (
                        <Button onClick={hook.handleSubmit}>
                            Submit
                        </Button>
                    )
                }
                </div>

        </>
    )
}

function StepForm({ hook }: { hook: ReturnType<typeof useCreateTemplate> }) {
    switch (hook.step) {
        case 0: return <StepBasicInfo  form={hook.form} />
        case 1: return <StepVariables variables={hook.variables} form={hook.form} />
        case 2: return <StepAddons addons={hook.addons} form={hook.form}/>
        case 3: return <StepFiles files={hook.files} form={hook.form} />
        case 4: return <StepReview values={hook.form.getValues()} />
        default: return null
    }
}