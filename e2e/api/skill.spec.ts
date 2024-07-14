import { test, expect } from "@playwright/test";
import { log } from "console";
import exp from "constants";

test("should response one skill when request POST /api/v1/skills", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	expect(reps.ok).toBeTruthy();
	const response = await reps.json();
	expect(response).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining({
				Key: "python",
				Name: "Python",
				Description:
					"Python is an interpreted, high-level, general-purpose programming language.",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
				Tags: ["programming language", "scripting"],
			}),
		})
	);

	const key = response["data"]["Key"];
	await request.delete("/api/v1/skills/" + key);
});

test("should response with error when request POST /api/v1/skills with duplicate key", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python 2",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});

	expect(reps.ok).toBeTruthy();
	const response = await reps.json();
	expect(response).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining({
				Key: "python 2",
				Name: "Python",
				Description:
					"Python is an interpreted, high-level, general-purpose programming language.",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
				Tags: ["programming language", "scripting"],
			}),
		})
	);

	const checkReps = await request.post("/api/v1/skills", {
		data: {
			key: "python 2",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});

	expect(!checkReps.ok()).toBeTruthy();
	const checkResponse = await checkReps.json();
	expect(checkResponse).toEqual(
		expect.objectContaining({
			status: "error",
			message: "Skill already exists",
		})
	);

	const key = response["data"]["Key"];
	await request.delete("/api/v1/skills/" + key);
});

test("should response with one skill when request GET /api/v1/skills/python3", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python3",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});

	const expectedResponse = {
		Key: "python3",
		Name: "Python",
		Description:
			"Python is an interpreted, high-level, general-purpose programming language.",
		Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
		Tags: ["programming language", "scripting"],
	};

	expect(reps.ok).toBeTruthy();
	const response = await reps.json();
	expect(response).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining(expectedResponse),
		})
	);

	const getReps = await request.get("/api/v1/skills/python3");
	expect(getReps.ok()).toBeTruthy();
	const getResponse = await getReps.json();
	expect(getResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining(expectedResponse),
		})
	);

	const key = response["data"]["Key"];
	await request.delete("/api/v1/skills/" + key);
});

test("should response with err when request GET /api/v1/skills/{key} when key is not in database", async ({
	request,
}) => {
	const getReps = await request.get("/api/v1/skills/python55");
	expect(getReps.ok()).toBeFalsy();
	const getResponse = await getReps.json();
	expect(getResponse).toEqual(
		expect.objectContaining({
			status: "error",
			message: "Skill not found",
		})
	);
});

test("should response with all skill when request GET /api/v1/skills", async ({
	request,
}) => {
	const reps1 = await request.post("/api/v1/skills", {
		data: {
			key: "pythontest1",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	const reps2 = await request.post("/api/v1/skills", {
		data: {
			key: "pythontest2",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	const reps3 = await request.post("/api/v1/skills", {
		data: {
			key: "pythontest3",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});

	expect(reps1.ok()).toBeTruthy();
	expect(reps2.ok()).toBeTruthy();
	expect(reps3.ok()).toBeTruthy();

	const getReps = await request.get("/api/v1/skills");
	expect(getReps.ok()).toBeTruthy();
	const getResponse = await getReps.json();
	expect(getResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.arrayContaining([
				{
					Key: "pythontest1",
					Name: "Python",
					Description:
						"Python is an interpreted, high-level, general-purpose programming language.",
					Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
					Tags: ["programming language", "scripting"],
				},
				{
					Key: "pythontest2",
					Name: "Python",
					Description:
						"Python is an interpreted, high-level, general-purpose programming language.",
					Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
					Tags: ["programming language", "scripting"],
				},
				{
					Key: "pythontest3",
					Name: "Python",
					Description:
						"Python is an interpreted, high-level, general-purpose programming language.",
					Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
					Tags: ["programming language", "scripting"],
				},
			]),
		})
	);

	await request.delete("/api/v1/skills/" + "pythontest1");
	await request.delete("/api/v1/skills/" + "pythontest2");
	await request.delete("/api/v1/skills/" + "pythontest3");
});

test("should response with nothing when request GET /api/v1/skills", async ({
	request,
}) => {
	const getReps = await request.get("/api/v1/skills");
	expect(getReps.ok()).toBeTruthy();
	const getResponse = await getReps.json();
	expect(getResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.arrayContaining([]),
		})
	);
});

test("should response with updated skill when PUT /api/v1/skills/{key} when key is available", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python10",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	expect(reps.ok).toBeTruthy();

	const updatedReps = await request.put("/api/v1/skills/python10", {
		data: {
			name: "Python 3",
			description:
				"Python 3 is the latest version of Python programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["data"],
		},
	});
  expect(updatedReps.ok()).toBeTruthy()
	const updateResponse = await updatedReps.json();
	expect(updateResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining({
				Key: "python10",
				Name: "Python 3",
				Description:
					"Python 3 is the latest version of Python programming language.",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
				Tags: ["data"],
			}),
		})
	);

  await request.delete("/api/v1/skills/python")
});

test("should response with error when PUT /api/v1/skills/{key} when key is unavailable", async ({
	request,
}) => {

	const updatedReps = await request.put("/api/v1/skills/python19", {
		data: {
			name: "Python 3",
			description:
				"Python 3 is the latest version of Python programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["data"],
		},
	});

  expect(updatedReps.ok()).toBeFalsy()
	const updateResponse = await updatedReps.json();
	expect(updateResponse).toEqual(
		expect.objectContaining({
			status: "error",
			message:"not be able to update skill"
		})
	);

  await request.delete("/api/v1/skills/python")
});

test("should response with updated skill name when PATCH /api/v1/skills/{key}/actions/nam when key is available", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python11",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	expect(reps.ok).toBeTruthy();

	const updatedReps = await request.patch("/api/v1/skills/python11/actions/name", {
		data: {
			name: "Python 3",
		},
	});
  expect(updatedReps.ok()).toBeTruthy()
	const updateResponse = await updatedReps.json();
	expect(updateResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining({
				Key: "python11",
				Name: "Python 3",
				Description:
					"Python is an interpreted, high-level, general-purpose programming language.",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
				Tags: ["programming language", "scripting"],
			}),
		})
	);

  await request.delete("/api/v1/skills/python")
});

test("should response with error when PATCH /api/v1/skills/{key}/actions/name when key is unavailable", async ({
	request,
}) => {

	const updatedReps = await request.patch("/api/v1/skills/python19/actions/name", {
		data: {
			name: "Python 3",
		},
	});

  // console.log(updatedReps)

  expect(updatedReps.ok()).toBeFalsy()
	const updateResponse = await updatedReps.json();
  
	expect(updateResponse).toEqual(
		expect.objectContaining({
			status: "error",
			message:"not be able to update skill name"
		})
	);
});

test("should response with updated skill when PATCH /api/v1/skills/{key}/actions/description when key is available", async ({
	request,
}) => {
	const reps = await request.post("/api/v1/skills", {
		data: {
			key: "python13",
			name: "Python",
			description:
				"Python is an interpreted, high-level, general-purpose programming language.",
			logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
			tags: ["programming language", "scripting"],
		},
	});
	expect(reps.ok).toBeTruthy();

	const updatedReps = await request.patch("/api/v1/skills/python13/actions/description", {
		data: {
			description:
				"Python 3 is the latest version of Python programming language.",
		},
	});
  expect(updatedReps.ok()).toBeTruthy()
	const updateResponse = await updatedReps.json();
	expect(updateResponse).toEqual(
		expect.objectContaining({
			status: "success",
			data: expect.objectContaining({
				Key: "python13",
				Name: "Python",
				Description:
					"Python 3 is the latest version of Python programming language.",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
				Tags: ["programming language", "scripting"],
			}),
		})
	);
})

//   await request.delete("/api/v1/skills/python")
// });

// test("should response with error when PATCH /api/v1/skills/{key}/actions/description when key is unavailable", async ({
// 	request,
// }) => {

// 	const updatedReps = await request.patch("/api/v1/skills/python19", {
// 		data: {
// 			description:
// 				"Python 3 is the latest version of Python programming language.",
// 		},
// 	});

//   expect(updatedReps.ok()).toBeFalsy()
// 	const updateResponse = await updatedReps.json();
// 	expect(updateResponse).toEqual(
// 		expect.objectContaining({
// 			status: "error",
// 			message:"not be able to update skill description"
// 		})
// 	);

//   await request.delete("/api/v1/skills/python")
// });
