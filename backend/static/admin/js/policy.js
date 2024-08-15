// URL to receive requests on server
const server_port = "8080";
const server_address = "http://localhost:" + server_port;
const policyEditUrl = server_address + "/admin/policy/";

// Function to edit policy from form submission
async function editPolicy(e, actionToComplete, role, resource, action) {
  //   Prevent default submit behavior
  e.preventDefault();

  //   Depending on action to complete, send request to server
  switch (actionToComplete) {
    case "create":
      // Send Post request to add policy
      response = await sendRequest("post", {
        role: role,
        resource: resource,
        action: action,
      });
      break;

    case "delete":
      response = await sendRequest("delete", {
        role: role,
        resource: resource,
        action: action,
      });
      break;
    case "role":
      getDetailFromRoleAddForm("role-add-form");

      response = await sendRequest("post", {
        role: role,
        resource: resource,
        action: action,
      });
      break;
    default:
      response = false;
  }

  if (!response) {
    alert("Action failed.");
  } else {
    // If successful
    // Once complete, Reload the page
    window.location.reload();
  }
}

// Used to add a role to a policy on a policy page
async function addRole(e, resource) {
  //   Prevent default submit behavior
  e.preventDefault();
  //   Get the form details
  const [role, action] = getDetailFromRoleAddForm("role-add-form");

  //   Send Post request to add policy
  response = await sendRequest("post", {
    role: role,
    resource: resource,
    action: action,
  });

  if (!response) {
    alert("Action failed.");
  } else {
    // If successful
    // Once complete, Reload the page
    window.location.reload();
  }
}

// Takes requested action in form: "POST", "DELETE" and sends requested action to server
// Takes data in form of JSON {role: role, resource: resource, action: action} and sends requested action to server
async function sendRequest(requestedAction, policyData) {
  try {
    // Convert requested action (POST/DELETE) to all caps
    capitalizedAction = requestedAction.toUpperCase();
    //   Build a param by replacing "/" with "-"
    slugParam = policyData.resource.replace(/\//g, "-");

    //  Build request URL
    requestUrl = policyEditUrl + slugParam;

    // Convert the data to JSON string
    policyDataJson = JSON.stringify(policyData);

    // Send a DELETE request to the server
    const response = await fetch(requestUrl, {
      method: capitalizedAction,
      headers: {
        "Content-Type": "application/json",
      },
      // Convert the selectedItems array to JSON and send it in the body of the request
      body: policyDataJson,
    });

    // If the response is not ok, then throw an error
    if (!response.ok) {
      return false;
    }

    // Reload the page
    return true;
  } catch (error) {
    console.error("Error:", error);
    return false;
  }
}

// Function to extract detail from form selectors
function getDetailFromRoleAddForm(formId) {
  //   Get the form element
  const form = document.getElementById(formId);

  //   Get the role and action from the form
  const role = form.role.value;
  const action = form.action.value;

  //   Return the role and action
  return [role, action];
}
