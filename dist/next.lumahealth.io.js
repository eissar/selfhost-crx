(() => {
  // src/globals.js
  var MRN_PARAM = "dt.mrn";

  // src/next.lumahealth.io/handle_parameters.js
  async function getIdForMRN(MRN) {
    const matched_patient = await fetch(
      "https://api.lumahealth.io/api/v2/patients?_select=name,dateOfBirth,contact,__t,medicalRecordNumber&_sort=name&patientGlobalSearchLimit=50&limit=10&_populate=chatActivity.assignee._id,chatActivity.assignee.name,chatActivity.assignee.type,chatActivity.assignee.roles,chatActivity.assignee.isOnline,chatActivity.groupAssignee.name,chatActivity.groupAssignee.type,chatActivity.automated,chatActivity.status&patientGlobalSearch=" + MRN,
      {
        headers: {
          accept: "application/json, text/plain, */*",
          "x-access-token": JSON.parse(
            localStorage.getItem("lh:auth_user")
          )?.token
        },
        referrer: "https://next.lumahealth.io/",
        body: null,
        method: "GET",
        mode: "cors",
        credentials: "omit"
      }
    ).then((r) => r.json()).then((j) => j.response[0]);
    if (!matched_patient) {
      alert("no match for MRN.");
    }
    return matched_patient._id;
  }
  async function HandleUrlParameters(url, params, mrn) {
    if (mrn) {
      params.delete(MRN_PARAM);
      const newUrl = url.origin + url.pathname + (params.toString() ? "?" + params.toString() : "") + url.hash;
      globalThis.history.replaceState({}, document.title, newUrl);
      const lumaId = await getIdForMRN(mrn);
      globalThis.location = `https://next.lumahealth.io/patients/${lumaId}/chat`;
    }
  }

  // src/next.lumahealth.io/main.js
  {
  }
  (async function() {
    const url = new URL(globalThis.location.href);
    const params = new URLSearchParams(url.search);
    const mrn = params.get(MRN_PARAM);
    if (mrn) {
      await HandleUrlParameters(url, params, mrn);
    }
  })();
})();
//# sourceMappingURL=next.lumahealth.io.js.map
